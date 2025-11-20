package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

func (v *Validator) Validate(data interface{}) error {
	if err := v.validate.Struct(data); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return formatValidationErrors(validationErrors)
		}
		return err
	}
	return nil
}

func formatValidationErrors(errs validator.ValidationErrors) error {
	var messages []string
	for _, err := range errs {
		message := fmt.Sprintf("field '%s' failed validation '%s'", err.Field(), err.Tag())
		if err.Param() != "" {
			message += fmt.Sprintf(" (param: %s)", err.Param())
		}
		messages = append(messages, message)
	}
	return fmt.Errorf("validation errors: %s", strings.Join(messages, "; "))
}

func (v *Validator) RegisterCustomValidation(tag string, fn validator.Func) error {
	return v.validate.RegisterValidation(tag, fn)
}
