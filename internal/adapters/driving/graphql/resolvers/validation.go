package resolvers

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// --- 1. Validator Setup ---

// Create a single, package-level validator instance for reuse.
// This is now accessible to all files in the 'resolvers' package.
var validate = validator.New()

// init() is a special Go function that runs once when the package is initialized.
// We use it to register our custom validation rules.
func init() {
	// Register the custom validation function for our 'emaildomain' tag.
	validate.RegisterValidation("emaildomain", validateEmailDomain)
	// Register the custom validation function for our 'passwordcomplex' tag.
	validate.RegisterValidation("passwordcomplex", validatePasswordComplexity)
}

// --- 2. Custom Validation Structs ---
// These structs mirror the GraphQL inputs but add validation tags.

// registerInputValidation maps to the RegisterInput, adding validation tags.
type registerInputValidation struct {
	Name     string `validate:"required,min=3,max=100"`
	Email    string `validate:"required,email,emaildomain"`
	Password string `validate:"required,min=8,max=32,passwordcomplex"`
}

// loginInputValidation maps to the LoginInput, adding validation tags.
type loginInputValidation struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

// tokenValidation is a reusable struct for validating the simple
// string input for RefreshToken, Logout, and VerifyEmail.
type tokenValidation struct {
	Token string `validate:"required"`
}

// resetPasswordInputValidation maps to the ResetPasswordInput.
type resetPasswordInputValidation struct {
	Token       string `validate:"required"`
	NewPassword string `validate:"required,min=8,max=32,passwordcomplex"`
}

// requestEmailChangeInputValidation maps to the RequestEmailChangeInput.
type requestEmailChangeInputValidation struct {
	NewEmail        string `validate:"required,email,emaildomain"`
	CurrentPassword string `validate:"required"`
}

// updateProfileInputValidation maps to the UpdateProfileInput.
// This validates the input based on requirement 1.11.2.
type updateProfileInputValidation struct {
	Name string `validate:"required,min=3,max=100"`
}

// disableAccountInputValidation maps to the DisableAccountInput.
// We just need to ensure the password field is present.
type disableAccountInputValidation struct {
	CurrentPassword string `validate:"required"`
}

// requestAccountDeletionInputValidation maps to RequestAccountDeletionInput.
type requestAccountDeletionInputValidation struct {
	CurrentPassword string `validate:"required"`
}

// --- 3. Custom Validation Functions ---

// validateEmailDomain implements the custom 'emaildomain' validation rule.
func validateEmailDomain(fl validator.FieldLevel) bool {
	// Get the email string from the field.
	email := fl.Field().String()

	// Check if the email ends with one of the allowed domains.
	if strings.HasSuffix(email, "@gmail.com") || strings.HasSuffix(email, "@uhb.ac.id") {
		return true
	}
	// If it doesn't match, validation fails.
	return false
}

// Pre-compile the regular expressions for password validation for efficiency.
var (
	hasUpper  = regexp.MustCompile(`[A-Z]`)
	hasLower  = regexp.MustCompile(`[a-z]`)
	hasNumber = regexp.MustCompile(`[0-9]`)
	hasSymbol = regexp.MustCompile(`[\W_]`) // \W is non-word, _ is often excluded by \W
)

// validatePasswordComplexity implements the 'passwordcomplex' validation rule.
func validatePasswordComplexity(fl validator.FieldLevel) bool {
	// Get the password string from the field.
	pass := fl.Field().String()

	// Check all required complexity rules from the PDF.
	if !hasUpper.MatchString(pass) {
		return false // Must contain one uppercase letter.
	}
	if !hasLower.MatchString(pass) {
		return false // Must contain one lowercase letter.
	}
	if !hasNumber.MatchString(pass) {
		return false // Must contain one number.
	}
	if !hasSymbol.MatchString(pass) {
		return false // Must contain one symbol.
	}
	// If all checks pass:
	return true
}

// --- 4. Error Translation Helper ---

// formatValidationErrors translates raw validation errors into simple, user-friendly messages.
func formatValidationErrors(err error) error {
	// This will hold our clean error message.
	var message string

	// We check if the error is the specific type from the validator library.
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		// We only care about the *first* error for simplicity.
		fieldError := validationErrors[0]

		// Get the field name (e.g., "Email") and the validation tag (e.g., "required").
		field := fieldError.Field()
		tag := fieldError.Tag()

		// Create user-friendly messages based on the tag.
		switch tag {
		case "required":
			message = fmt.Sprintf("%s is required", field)
		case "email":
			message = fmt.Sprintf("%s is not a valid email address", field)
		case "min":
			message = fmt.Sprintf("%s must be at least %s characters long", field, fieldError.Param())
		case "max":
			message = fmt.Sprintf("%s must not be more than %s characters long", field, fieldError.Param())
		case "emaildomain": // Our custom tag
			message = "Email must be a @gmail.com or @uhb.ac.id address"
		case "passwordcomplex": // Our custom tag
			message = "Password must contain an uppercase letter, a lowercase letter, a number, and a symbol"
		default:
			// A fallback for any other validation rules.
			message = fmt.Sprintf("Field %s failed validation on tag '%s'", field, tag)
		}
	} else {
		// This is not a validation error, so just return the original error.
		return err
	}

	// Return a new, simple error with our clean message.
	return errors.New(message)
}
