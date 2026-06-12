package handler

import (
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/Pavan20008/user-age-api/internal/models"
)

// newValidator builds a validator with a custom "dob" rule that ensures the
// value is a valid ISO date (YYYY-MM-DD) and is not in the future.
func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	_ = v.RegisterValidation("dob", validateDOB)
	return v
}

func validateDOB(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	t, err := time.Parse(models.DateLayout, value)
	if err != nil {
		return false
	}
	// Date of birth cannot be in the future.
	return !t.After(time.Now().UTC())
}
