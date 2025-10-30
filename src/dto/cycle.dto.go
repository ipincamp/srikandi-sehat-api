package dto

import (
	"time"
)

// --- Generic Pagination Structures ---
type Pagination struct {
	Limit        int   `json:"limit"`
	TotalRows    int64 `json:"total_data"`
	TotalPages   int   `json:"total_pages"`
	CurrentPage  int   `json:"current_page"`
	PreviousPage *int  `json:"previous_page,omitempty"`
	NextPage     *int  `json:"next_page,omitempty"`
}

type PaginatedResponse[T any] struct {
	Data     []T        `json:"data"`
	Metadata Pagination `json:"metadata"`
}

// --- Cycle Specific ---

// CycleExportRecord defines the flattened structure for the cycle data CSV export.
type CycleExportRecord struct {
	// UserUUID       string `json:"user_uuid"`
	UserName       string `json:"user_name"`
	CycleNumber    int64  `json:"cycle_number"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	PeriodLength   int16  `json:"period_length"`
	PeriodCategory string `json:"period_category"`
	CycleLength    int16  `json:"cycle_length"`
	CycleCategory  string `json:"cycle_category"`
	Symptoms       string `json:"symptoms"`
}

// Request Param
type CycleParam struct {
	ID uint `params:"id" validate:"required,numeric"`
}

// Request Query
type PaginationQuery struct {
	Page  int `query:"page" validate:"omitempty,numeric,min=1"`
	Limit int `query:"limit" validate:"omitempty,numeric,min=1"`
}

// Request Body
type CycleRequest struct {
	StartDate string `json:"start_date" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	EndDate   string `json:"finish_date" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

type DeleteCycleRequest struct {
	Reason string `json:"reason" validate:"required,min=5,max=255"`
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

type CycleStatusResponse struct {
	IsOnCycle           bool    `json:"is_on_cycle"`
	CurrentPeriodDay    *int    `json:"current_period_day,omitempty"`
	IsPeriodNormal      *bool   `json:"is_period_normal,omitempty"`
	LastPeriodLength    *int    `json:"last_period_length,omitempty"`
	CurrentCycleLength  *int    `json:"current_cycle_length,omitempty"`
	LastCycleLength     *int    `json:"last_cycle_length,omitempty"`
	IsCycleNormal       *bool   `json:"is_cycle_normal,omitempty"`
	DaysUntilNextPeriod *int    `json:"days_until_next_period,omitempty"`
	PredictedPeriodDate *string `json:"predicted_period_date,omitempty"`
	Message             string  `json:"message"`
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
