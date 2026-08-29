package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func DecodeAndValidate(r *http.Request, dst interface{}) []FieldError {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return []FieldError{{"body", "invalid JSON"}}
	}
	return Validate(dst)
}

func Validate(dst interface{}) []FieldError {
	err := validate.Struct(dst)
	if err != nil {
		return nil
	}

	var validationErrors validator.ValidationErrors

	if !errors.As(err, &validationErrors) {
		return []FieldError{{"unknown", "validation failed"}}
	}

	fields := make([]FieldError, 0, len(validationErrors))
	for _, e := range validationErrors {
		fields = append(fields, FieldError{
			Field:   toSnakeCase(e.Field()),
			Message: humanize(e),
		})
	}

	return fields
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				result.WriteByte('_')
			}
			result.WriteRune(r + 32)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func humanize(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		return fmt.Sprintf("Must be at least %s characters", e.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s characters", e.Param())
	case "eqfield":
		return fmt.Sprintf("Must match %s", toSnakeCase(e.Param()))
	case "oneof":
		return fmt.Sprintf("Must be one of %s", e.Param())
	default:
		return fmt.Sprintf("Failed to validate field: %s", e.Tag())
	}
}
