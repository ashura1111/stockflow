package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

func FormatValidationError(err error) string {
	if err == nil {
		return ""
	}
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			return fmt.Sprintf("Validation failed on field '%s': %s", fieldError.Field(), fieldError.Tag())
		}
	}
	return err.Error()
}
