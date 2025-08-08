package dto

import "time"

// Request Param
type CycleParam struct {
	ID uint `params:"id" validate:"required,numeric"`
}

// Request Body
type CycleRequest struct {
	StartDate string `json:"start_date" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	EndDate   string `json:"finish_date" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

// Response Body
type CycleResponse struct {
	ID             uint       `json:"id"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        *time.Time `json:"finish_date,omitempty"`
	PeriodLength   *int16     `json:"period_length,omitempty"`
	CycleLength    *int16     `json:"cycle_length,omitempty"`
	IsPeriodNormal *bool      `json:"is_period_normal,omitempty"`
	IsCycleNormal  *bool      `json:"is_cycle_normal,omitempty"`
}

type SymptomLogInCycleResponse struct {
	LoggedAt        time.Time `json:"logged_at"`
	Note            *string   `json:"note,omitempty"`
	SymptomName     string    `json:"symptom_name"`
	SymptomCategory string    `json:"symptom_category"`
	SelectedOption  *string   `json:"selected_option,omitempty"`
}

type SymptomDetail struct {
	SymptomName     string  `json:"symptom_name"`
	SymptomCategory string  `json:"symptom_category"`
	SelectedOption  *string `json:"selected_option,omitempty"`
}

type SymptomLogGroupResponse struct {
	ID       uint            `json:"id"`
	LoggedAt time.Time       `json:"logged_at"`
	Note     *string         `json:"note,omitempty"`
	Details  []SymptomDetail `json:"details"`
}

type CycleDetailResponse struct {
	ID             uint                      `json:"id"`
	StartDate      time.Time                 `json:"start_date"`
	EndDate        *time.Time                `json:"finish_date,omitempty"`
	PeriodLength   *int16                    `json:"period_length,omitempty"`
	CycleLength    *int16                    `json:"cycle_length,omitempty"`
	IsPeriodNormal *bool                     `json:"is_period_normal,omitempty"`
	IsCycleNormal  *bool                     `json:"is_cycle_normal,omitempty"`
	Symptoms       []SymptomLogGroupResponse `json:"symptoms"`
}
