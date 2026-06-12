package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/Pavan20008/user-age-api/config"
	"github.com/Pavan20008/user-age-api/db"
	"github.com/Pavan20008/user-age-api/internal/handler"
	"github.com/Pavan20008/user-age-api/internal/logger"
	"github.com/Pavan20008/user-age-api/internal/middleware"
	"github.com/Pavan20008/user-age-api/internal/repository"
	"github.com/Pavan20008/user-age-api/internal/routes"
	"github.com/Pavan20008/user-age-api/internal/service"
)

func main() {
	cfg := config.Load()

	log, err := logger.New(cfg.Environment)
	if err != nil {
		panic(err)
	}
	defer func() { _ = log.Sync() }()

	// Connect to PostgreSQL.
	ctx := context.Background()
	pool, err := connectDB(ctx, cfg.DatabaseURL, log)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	// Apply migrations on startup.
	if err := db.Migrate(ctx, pool); err != nil {
		log.Fatal("failed to run migrations", zap.Error(err))
	}
	log.Info("migrations applied")

	// Wire the layers together.
	repo := repository.NewUserRepository(pool)
	svc := service.NewUserService(repo, log)
	h := handler.NewUserHandler(svc, log)

	app := fiber.New(fiber.Config{
		ErrorHandler:          handler.ErrorHandler,
		DisableStartupMessage: true,
	})
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger(log))

	routes.Register(app, h)

	// Start the server and handle graceful shutdown.
	go func() {
		addr := ":" + cfg.AppPort
		log.Info("starting server", zap.String("addr", addr))
		if err := app.Listen(addr); err != nil {
			log.Fatal("server stopped", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server")
	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		log.Error("graceful shutdown failed", zap.Error(err))
	}
}

func connectDB(ctx context.Context, url string, log *zap.Logger) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}

	// Retry the initial ping so the app can start alongside a freshly booted
	// database (e.g. in docker-compose).
	var pingErr error
	for i := 0; i < 10; i++ {
		pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		pingErr = pool.Ping(pingCtx)
		cancel()
		if pingErr == nil {
			return pool, nil
		}
		log.Warn("waiting for database", zap.Int("attempt", i+1), zap.Error(pingErr))
		time.Sleep(2 * time.Second)
	}
	pool.Close()
	return nil, errors.New("database not reachable: " + pingErr.Error())
}
