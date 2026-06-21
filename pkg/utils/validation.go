package utils

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ExtractValidationErrors(err error) []FieldError {
	var fieldErrors []FieldError
	if !errors.As(err, &validator.ValidationErrors{}) {
		return nil
	}
	// Assuming err is of type validator.ValidationErrors
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, ve := range validationErrors {
			fieldErrors = append(fieldErrors, FieldError{
				Field:   ve.Field(),
				Message: fieldErrorMessage(ve),
			})
		}
	}
	return fieldErrors
}

func fieldErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "email":
		return fe.Field() + " has an invalid email format"
	case "min":
		return fe.Field() + " is too short. Minimum length is " + fe.Param() + " characters"
	case "max":
		return fe.Field() + " is too long. Maximum length is " + fe.Param() + " characters"
	default:
		return fe.Field() + " is invalid" // Default error message
	}
}
