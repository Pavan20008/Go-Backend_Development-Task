package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RequestIDHeader is the header used to carry the request correlation id.
const RequestIDHeader = "X-Request-ID"

// RequestID ensures every request has a correlation id. It reuses an incoming
// X-Request-ID header when present, otherwise it generates one, stores it in
// locals, and echoes it back on the response.
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Get(RequestIDHeader)
		if id == "" {
			id = uuid.NewString()
		}
		c.Locals("requestid", id)
		c.Set(RequestIDHeader, id)
		return c.Next()
	}
}

// RequestLogger logs the method, path, status and duration of each request
// using Uber Zap, tagging entries with the request id.
func RequestLogger(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)
		requestID, _ := c.Locals("requestid").(string)

		logger.Info("request",
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", c.OriginalURL()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("duration", duration),
			zap.String("ip", c.IP()),
		)

		return err
	}
}
