package dto

import (
	"ipincamp/srikandi-sehat/src/models"
	"time"
)

// Request Body
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Response Body
type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role,omitempty"`
	Token     string    `json:"token,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func AuthResponseJson(user models.User, token ...string) UserResponse {
	var roleName string
	if len(user.Roles) > 0 {
		roleName = user.Roles[0].Name
	}

	response := UserResponse{
		ID:        user.UUID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      roleName,
		CreatedAt: user.CreatedAt,
	}

	if len(token) > 0 {
		response.Token = token[0]
	}

	return response
}
