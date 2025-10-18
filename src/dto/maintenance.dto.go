package dto

type MaintenanceStatusResponse struct {
	IsMaintenance bool   `json:"active"`
	Message       string `json:"message,omitempty"`
}

type ToggleMaintenanceRequest struct {
	Active  *bool  `json:"active" validate:"required"`
	Message string `json:"message" validate:"omitempty"`
}

type WhitelistUserRequest struct {
	UserUUID string `json:"user_uuid" validate:"required,uuid"`
}

type WhitelistedUserResponse struct {
	UserUUID string `json:"user_uuid"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}
