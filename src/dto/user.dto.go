package dto

import (
	"fmt"
	"ipincamp/srikandi-sehat/src/constants"
	"ipincamp/srikandi-sehat/src/models"
	"math"
	"strings"
	"time"
)

// --- Request Params & Queries ---
type UserParam struct {
	ID string `params:"id" validate:"required,uuid"`
}

type UserQuery struct {
	Classification string `query:"classification" validate:"omitempty,oneof=urban rural"`
	Page           int    `query:"page" validate:"omitempty,numeric,min=1"`
	Limit          int    `query:"limit" validate:"omitempty,numeric,min=1,max=100"`
}

// --- Request Body ---
type UpdateProfileRequest struct {
	Name                *string                   `json:"name" validate:"omitempty,min=3"`
	PhoneNumber         *string                   `json:"phone" validate:"omitempty,min=10,max=15"`
	VillageCode         *string                   `json:"address_code" validate:"omitempty,len=10"`
	DateOfBirth         *string                   `json:"birthdate" validate:"omitempty,datetime=2006-01-02"`
	HeightCM            *uint                     `json:"tb_cm" validate:"omitempty,gte=100,lte=250"`
	WeightKG            *float32                  `json:"bb_kg" validate:"omitempty,gte=30,lte=200"`
	LastEducation       *constants.EducationLevel `json:"edu_now" validate:"omitempty,oneof='Tidak Sekolah' SD SMP SMA Diploma S1 S2 S3"`
	ParentLastEducation *constants.EducationLevel `json:"edu_parent" validate:"omitempty,oneof='Tidak Sekolah' SD SMP SMA Diploma S1 S2 S3"`
	ParentLastJob       *string                   `json:"job_parent" validate:"omitempty"`
	InternetAccess      *constants.InternetAccess `json:"inet_access" validate:"omitempty,oneof=WiFi Seluler"`
	MenarcheAge         *uint                     `json:"first_haid" validate:"omitempty,gte=8,lte=20"`
}

type ChangePasswordRequest struct {
	OldPassword             string `json:"old_password" validate:"required"`
	NewPassword             string `json:"new_password" validate:"required,min=8,password_strength"`
	NewPasswordConfirmation string `json:"new_password_confirmation" validate:"required,eqfield=NewPassword"`
}

// --- Response Body ---
type CycleHistoryEntry struct {
	ID               uint       `json:"id"`
	StartDate        time.Time  `json:"start_date"`
	FinishDate       *time.Time `json:"finish_date,omitempty"`
	PeriodLengthDays *int16     `json:"period_length_days"`
	CycleLengthDays  *int16     `json:"cycle_length_days"`
	DeletedAt        *time.Time `json:"deleted_at"`
	DeletionReason   *string    `json:"deletion_reason"`
}

type ProfileResponse struct {
	PhoneNumber         string                   `json:"phone"`
	DateOfBirth         *time.Time               `json:"birthdate"`
	HeightCM            uint                     `json:"tb_cm"`
	WeightKG            float32                  `json:"bb_kg"`
	Bmi                 float32                  `json:"bmi,omitempty"`
	LastEducation       constants.EducationLevel `json:"edu_now"`
	ParentLastEducation constants.EducationLevel `json:"edu_parent"`
	ParentLastJob       string                   `json:"job_parent"`
	InternetAccess      constants.InternetAccess `json:"inet_access"`
	MenarcheAge         uint                     `json:"first_haid"`
	Address             string                   `json:"address"`
	UpdatedAt           *time.Time               `json:"updated_at,omitempty"`
}

type UserResponse struct {
	ID                string              `json:"id"`
	Name              string              `json:"name"`
	Email             string              `json:"email"`
	Role              string              `json:"role,omitempty"`
	Token             string              `json:"token,omitempty"`
	IsProfileComplete bool                `json:"profile_complete"`
	Profile           *ProfileResponse    `json:"profile,omitempty"`
	CycleHistory      []CycleHistoryEntry `json:"cycle_history,omitempty"`
	CreatedAt         time.Time           `json:"created_at"`
}

type UserStatisticsResponse struct {
	TotalRuralUsers  int64 `json:"total_rural_users"`
	TotalUrbanUsers  int64 `json:"total_urban_users"`
	TotalActiveUsers int64 `json:"total_active_users"`
	TotalUsers       int64 `json:"total_users"`
}

type UserCSVRecord struct {
	UUID                string    `json:"uuid"`
	Name                string    `json:"name"`
	Email               string    `json:"email"`
	DateOfBirth         *string   `json:"date_of_birth"`
	PhoneNumber         string    `json:"phone_number"`
	HeightCM            uint      `json:"height_cm"`
	WeightKG            float32   `json:"weight_kg"`
	BMI                 float32   `json:"bmi"`
	BMICategory         string    `json:"bmi_category"`
	MenarcheAge         uint      `json:"menarche_age"`
	LastEducation       string    `json:"last_education"`
	ParentLastEducation string    `json:"parent_last_education"`
	ParentLastJob       string    `json:"parent_last_job"`
	InternetAccess      string    `json:"internet_access"`
	Village             string    `json:"village"`
	District            string    `json:"district"`
	Regency             string    `json:"regency"`
	Province            string    `json:"province"`
	Classification      string    `json:"classification"`
	RegisteredAt        time.Time `json:"registered_at"`
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
		address := buildFullAddress(user.Profile)

		var bmi float32
		if user.Profile.HeightCM > 0 && user.Profile.WeightKG > 0 {
			heightInMeters := float32(user.Profile.HeightCM) / 100
			bmi = user.Profile.WeightKG / (heightInMeters * heightInMeters)
			bmi = float32(math.Round(float64(bmi)*100) / 100)
		}

		profileData = &ProfileResponse{
			PhoneNumber:         user.Profile.PhoneNumber,
			DateOfBirth:         user.Profile.DateOfBirth,
			HeightCM:            user.Profile.HeightCM,
			WeightKG:            user.Profile.WeightKG,
			Bmi:                 bmi,
			LastEducation:       user.Profile.LastEducation,
			ParentLastEducation: user.Profile.ParentLastEducation,
			ParentLastJob:       user.Profile.ParentLastJob,
			InternetAccess:      user.Profile.InternetAccess,
			MenarcheAge:         user.Profile.MenarcheAge,
			Address:             address,
			UpdatedAt:           &user.Profile.UpdatedAt,
		}
	}

	var response UserResponse
	if len(token) > 0 {
		response = UserResponse{
			ID:                user.UUID,
			Name:              user.Name,
			Email:             user.Email,
			Role:              roleName,
			Token:             token[0],
			IsProfileComplete: isProfileComplete,
			CreatedAt:         user.CreatedAt,
		}
	} else {
		response = UserResponse{
			ID:                user.UUID,
			Name:              user.Name,
			Email:             user.Email,
			Role:              roleName,
			IsProfileComplete: isProfileComplete,
			Profile:           profileData,
			CreatedAt:         user.CreatedAt,
		}
	}

	return response
}

func buildFullAddress(profile models.Profile) string {
	addressParts := []string{}

	if profile.Village.ID > 0 {
		classification := profile.Village.Classification.Name
		if strings.EqualFold(classification, string(constants.RuralClassification)) {
			classification = "DESA"
		} else if strings.EqualFold(classification, string(constants.UrbanClassification)) {
			classification = "KOTA"
		}
		village := profile.Village.Name

		addressParts = append(addressParts, fmt.Sprintf("(%s) %s", strings.ToUpper(classification), strings.ToUpper(village)))

		if profile.Village.District.ID > 0 {
			district := profile.Village.District.Name
			addressParts = append(addressParts, fmt.Sprintf("KECAMATAN %s", strings.ToUpper(district)))

			if profile.Village.District.Regency.ID > 0 {
				regency := profile.Village.District.Regency.Name
				addressParts = append(addressParts, fmt.Sprintf("KABUPATEN %s", strings.ToUpper(regency)))

				if profile.Village.District.Regency.Province.ID > 0 {
					province := profile.Village.District.Regency.Province.Name
					addressParts = append(addressParts, fmt.Sprintf("PROVINSI %s", strings.ToUpper(province)))
				}
			}
		}
	}

	return strings.Join(addressParts, ", ")
}
