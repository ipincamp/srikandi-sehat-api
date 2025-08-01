package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func SetupValidator() {
	validate = validator.New()
	validate.RegisterValidation("password_strength", ValidatePasswordStrength)
}

func ValidatePasswordStrength(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	hasUpper, _ := regexp.MatchString(`[A-Z]`, password)
	hasLower, _ := regexp.MatchString(`[a-z]`, password)
	hasNumber, _ := regexp.MatchString(`[0-9]`, password)
	hasSymbol, _ := regexp.MatchString(`[\W_]`, password)
	return hasUpper && hasLower && hasNumber && hasSymbol
}

type ValidationError struct {
	FailedField string `json:"field"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
}

func ValidateStruct(payload interface{}) []*ValidationError {
	var errors []*ValidationError
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidationError
			element.FailedField = err.Field()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
