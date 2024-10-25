package validations

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Property string `json:"property"`
	Tag      string `json:"tag"`
	Value    string `json:"value"`
	Message  string `json:"message"`
}

func GetValidationErrors(err error) *[]ValidationError {
	var validationError []ValidationError
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidationError
			element.Property = err.Field()
			element.Tag = err.Tag()
			element.Value = err.Param()
			validationError = append(validationError, element)
		}
		return &validationError
	}
	return nil
}
