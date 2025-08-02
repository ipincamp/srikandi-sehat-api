package dto

import (
	"ipincamp/srikandi-sehat/src/constants"
	"ipincamp/srikandi-sehat/src/models"
	"time"
)

// Request Body
type UpdateProfileRequest struct {
	Name                string                   `json:"name" validate:"omitempty,min=3"`
	PhoneNumber         string                   `json:"phone_number" validate:"required,min=10,max=15"`
	DateOfBirth         string                   `json:"date_of_birth" validate:"required,datetime=2006-01-02"`
	HeightCM            uint                     `json:"height_cm" validate:"required,gte=100,lte=250"`
	WeightKG            float32                  `json:"weight_kg" validate:"required,gte=30,lte=200"`
	AddressStreet       string                   `json:"address_street" validate:"required"`
	VillageCode         string                   `json:"village_code" validate:"required,len=10"`
	LastEducation       constants.EducationLevel `json:"last_education" validate:"required,oneof='Tidak Sekolah' SD SMP SMA Diploma S1 S2 S3"`
	ParentLastEducation constants.EducationLevel `json:"parent_last_education" validate:"required,oneof='Tidak Sekolah' SD SMP SMA Diploma S1 S2 S3"`
	ParentLastJob       string                   `json:"parent_last_job" validate:"required"`
	InternetAccess      constants.InternetAccess `json:"internet_access" validate:"required,oneof=WiFi Seluler"`
	MenarcheAge         uint                     `json:"menarche_age" validate:"required,gte=8,lte=20"`
}

type ChangePasswordRequest struct {
	OldPassword             string `json:"old_password" validate:"required"`
	NewPassword             string `json:"new_password" validate:"required,min=8,password_strength"`
	NewPasswordConfirmation string `json:"new_password_confirmation" validate:"required,eqfield=NewPassword"`
}

// Response Body
type ProfileResponse struct {
	PhoneNumber string     `json:"phone_number"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	HeightCM    uint       `json:"height_cm"`
	WeightKG    float32    `json:"weight_kg"`
}

type UserResponse struct {
	ID                string           `json:"id"`
	Name              string           `json:"name"`
	Email             string           `json:"email"`
	Role              string           `json:"role,omitempty"`
	Token             string           `json:"token,omitempty"`
	IsProfileComplete bool             `json:"is_profile_complete"`
	Profile           *ProfileResponse `json:"profile,omitempty"`
	CreatedAt         time.Time        `json:"created_at"`
}

// Formatter
func UserResponseJson(user models.User, token ...string) UserResponse {
	var roleName string
	if len(user.Roles) > 0 {
		roleName = user.Roles[0].Name
	}

	isProfileComplete := user.Profile.ID > 0
	var profileData *ProfileResponse
	if isProfileComplete {
		profileData = &ProfileResponse{
			PhoneNumber: user.Profile.PhoneNumber,
			DateOfBirth: user.Profile.DateOfBirth,
			HeightCM:    user.Profile.HeightCM,
			WeightKG:    user.Profile.WeightKG,
		}
	}

	response := UserResponse{
		ID:                user.UUID,
		Name:              user.Name,
		Email:             user.Email,
		Role:              roleName,
		IsProfileComplete: isProfileComplete,
		Profile:           profileData,
		CreatedAt:         user.CreatedAt,
	}

	if len(token) > 0 {
		response.Token = token[0]
	}

	return response
}
