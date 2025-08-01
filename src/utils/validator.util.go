package utils

import (
	"regexp"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var validate *validator.Validate
var trans ut.Translator

func SetupValidator() {
	validate = validator.New()

	english := en.New()
	uni := ut.New(english, english)
	trans, _ = uni.GetTranslator("en")

	en_translations.RegisterDefaultTranslations(validate, trans)

	validate.RegisterValidation("password_strength", ValidatePasswordStrength)
	validate.RegisterTranslation("password_strength", trans, func(ut ut.Translator) error {
		return ut.Add("password_strength", "The {0} field must contain at least one uppercase letter, lowercase letter, number, and symbol.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("password_strength", fe.Field())
		return t
	})
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
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

func ValidateStruct(payload interface{}) []*ValidationError {
	var errors []*ValidationError
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidationError
			element.Field = err.Field()
			element.Tag = err.Tag()
			element.Value = err.Param()
			element.Message = err.Translate(trans)
			errors = append(errors, &element)
		}
	}
	return errors
}
