package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Pavan20008/user-age-api/internal/handler"
)

// Register wires the user routes onto the Fiber app.
func Register(app *fiber.App, h *handler.UserHandler) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	users := app.Group("/users")
	users.Post("/", h.Create)
	users.Get("/", h.List)
	users.Get("/:id", h.Get)
	users.Put("/:id", h.Update)
	users.Delete("/:id", h.Delete)
}
