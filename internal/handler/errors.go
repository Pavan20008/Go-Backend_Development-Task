package handler

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ValidationError represents a 400 response carrying per-field details.
type ValidationError struct {
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields"`
}

func (e *ValidationError) Error() string { return e.Message }

// validationError converts go-playground/validator output into a structured
// ValidationError that the global error handler renders as JSON.
func validationError(err error) error {
	var verrs validator.ValidationErrors
	fields := map[string]string{}
	if ok := asValidationErrors(err, &verrs); ok {
		for _, fe := range verrs {
			fields[fe.Field()] = messageForTag(fe)
		}
	}
	return &ValidationError{
		Message: "validation failed",
		Fields:  fields,
	}
}

func asValidationErrors(err error, target *validator.ValidationErrors) bool {
	if verrs, ok := err.(validator.ValidationErrors); ok {
		*target = verrs
		return true
	}
	return false
}

func messageForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "dob":
		return "must be a valid date (YYYY-MM-DD) and not in the future"
	case "min":
		return fmt.Sprintf("must be at least %s characters", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", fe.Param())
	default:
		return fmt.Sprintf("failed %q validation", fe.Tag())
	}
}

// ErrorHandler is the centralized Fiber error handler. It renders a consistent
// JSON error envelope and selects an appropriate HTTP status code.
func ErrorHandler(c *fiber.Ctx, err error) error {
	if ve, ok := err.(*ValidationError); ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  ve.Message,
			"fields": ve.Fields,
		})
	}

	code := fiber.StatusInternalServerError
	msg := "internal server error"
	if fe, ok := err.(*fiber.Error); ok {
		code = fe.Code
		msg = fe.Message
	}

	return c.Status(code).JSON(fiber.Map{"error": msg})
}
