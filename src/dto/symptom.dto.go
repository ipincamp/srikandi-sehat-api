package dto

import (
	"ipincamp/srikandi-sehat/src/constants"
	"time"
)

// Request Query
type SymptomLogQuery struct {
	StartDate string `query:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate   string `query:"finish_date" validate:"required,datetime=2006-01-02"`
}

type RecommendationQuery struct {
	SymptomIDs string `query:"symptom_ids" validate:"required"`
}

type SymptomHistoryQuery struct {
	Page      int    `query:"page" validate:"omitempty,numeric,min=1"`
	Limit     int    `query:"limit" validate:"omitempty,numeric,min=1"`
	StartDate string `query:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate   string `query:"finish_date" validate:"omitempty,datetime=2006-01-02"`
	Date      string `query:"date" validate:"omitempty,datetime=2006-01-02"`
}

// Request Body
type SymptomLogDetailRequest struct {
	SymptomID       uint  `json:"symptom_id" validate:"required"`
	SymptomOptionID *uint `json:"option_id,omitempty"`
}

type SymptomLogRequest struct {
	LoggedAt string                    `json:"logged_at" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Note     string                    `json:"note" validate:"omitempty"`
	Symptoms []SymptomLogDetailRequest `json:"symptoms" validate:"required,min=1"`
}

// Response Body
type SymptomOptionResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type SymptomMasterResponse struct {
	ID      uint                    `json:"id"`
	Name    string                  `json:"name"`
	Type    constants.SymptomType   `json:"type"`
	Options []SymptomOptionResponse `json:"options,omitempty"`
}

type SymptomLogDetailResponse struct {
	SymptomName     string `json:"symptom_name"`
	SymptomCategory string `json:"symptom_category"`
	SelectedOption  string `json:"selected_option,omitempty"`
}

type SymptomLogResponse struct {
	LoggedAt time.Time                  `json:"logged_at"`
	Note     string                     `json:"note,omitempty"`
	Details  []SymptomLogDetailResponse `json:"details"`
}

// SymptomHistoryResponse defines the structure for the aggregated symptom history list.
type SymptomHistoryResponse struct {
	ID            uint   `json:"id"`
	TotalSymptoms int    `json:"total_symptoms"`
	LogDate       string `json:"log_date"`
}

type RecommendationResponse struct {
	ForSymptom  string `json:"for_symptom"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Source      string `json:"source,omitempty"`
}
