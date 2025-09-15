package dto

// Request Body
type RegisterRequest struct {
	Name                 string `json:"name" validate:"required,min=3"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=8,password_strength"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
	FCMToken             string `json:"fcm_token" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
