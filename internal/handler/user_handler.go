package handler

import (
	"errors"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/Pavan20008/user-age-api/internal/models"
	"github.com/Pavan20008/user-age-api/internal/service"
)

const (
	defaultPage     = 1
	defaultPageSize = 10
	maxPageSize     = 100
)

// UserHandler exposes the HTTP handlers for user resources.
type UserHandler struct {
	svc      *service.UserService
	validate *validator.Validate
	logger   *zap.Logger
}

// NewUserHandler builds a UserHandler.
func NewUserHandler(svc *service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		svc:      svc,
		validate: newValidator(),
		logger:   logger,
	}
}

// Create handles POST /users.
func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON body")
	}
	if err := h.validate.Struct(req); err != nil {
		return validationError(err)
	}

	user, err := h.svc.Create(c.Context(), req)
	if err != nil {
		return mapServiceError(err)
	}
	return c.Status(fiber.StatusCreated).JSON(user)
}

// Get handles GET /users/:id.
func (h *UserHandler) Get(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}

	user, err := h.svc.Get(c.Context(), id)
	if err != nil {
		return mapServiceError(err)
	}
	return c.JSON(user)
}

// List handles GET /users with optional ?page= and ?limit= pagination.
func (h *UserHandler) List(c *fiber.Ctx) error {
	page := queryInt(c, "page", defaultPage)
	limit := queryInt(c, "limit", defaultPageSize)
	if page < 1 {
		page = defaultPage
	}
	if limit < 1 {
		limit = defaultPageSize
	}
	if limit > maxPageSize {
		limit = maxPageSize
	}
	offset := (page - 1) * limit

	users, err := h.svc.List(c.Context(), int32(limit), int32(offset))
	if err != nil {
		return mapServiceError(err)
	}

	total, err := h.svc.Count(c.Context())
	if err != nil {
		return mapServiceError(err)
	}

	// Pagination metadata is exposed via headers so the body stays a plain
	// JSON array as specified by the API contract.
	c.Set("X-Total-Count", strconv.FormatInt(total, 10))
	c.Set("X-Page", strconv.Itoa(page))
	c.Set("X-Limit", strconv.Itoa(limit))

	return c.JSON(users)
}

// Update handles PUT /users/:id.
func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON body")
	}
	if err := h.validate.Struct(req); err != nil {
		return validationError(err)
	}

	user, err := h.svc.Update(c.Context(), id, req)
	if err != nil {
		return mapServiceError(err)
	}
	return c.JSON(user)
}

// Delete handles DELETE /users/:id.
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}

	if err := h.svc.Delete(c.Context(), id); err != nil {
		return mapServiceError(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func parseID(c *fiber.Ctx) (int32, error) {
	raw := c.Params("id")
	id, err := strconv.ParseInt(raw, 10, 32)
	if err != nil || id <= 0 {
		return 0, fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}
	return int32(id), nil
}

func queryInt(c *fiber.Ctx, key string, fallback int) int {
	raw := c.Query(key)
	if raw == "" {
		return fallback
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return v
}

func mapServiceError(err error) error {
	if errors.Is(err, service.ErrNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}
	if errors.Is(err, service.ErrInvalidDate) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid date of birth, expected format YYYY-MM-DD")
	}
	return err
}
