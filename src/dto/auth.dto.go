package dto

// Request Body
type RegisterRequest struct {
	Name                 string `json:"name" validate:"required,min=3"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=8,password_strength"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type VerifyOTPRequest struct {
	OTP string `json:"otp" validate:"required,len=6,numeric"`
}
