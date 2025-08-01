package dto

// Request Body
type UpdateDetailsRequest struct {
	Name  string `json:"name" validate:"omitempty,min=3"`
	Email string `json:"email" validate:"omitempty,email"`
}

type ChangePasswordRequest struct {
	OldPassword             string `json:"old_password" validate:"required"`
	NewPassword             string `json:"new_password" validate:"required,min=8,password_strength"`
	NewPasswordConfirmation string `json:"new_password_confirmation" validate:"required,eqfield=NewPassword"`
}
