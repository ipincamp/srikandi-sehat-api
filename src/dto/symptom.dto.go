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

// Request Body
type SymptomLogDetailRequest struct {
	SymptomID       uint  `json:"symptom_id" validate:"required"`
	SymptomOptionID *uint `json:"option_id,omitempty"`
}

type SymptomLogRequest struct {
	LogDate  string                    `json:"log_date" validate:"required,datetime=2006-01-02"`
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
	LogDate time.Time                  `json:"log_date"`
	Note    string                     `json:"note,omitempty"`
	Details []SymptomLogDetailResponse `json:"details"`
}

type RecommendationResponse struct {
	ForSymptom  string `json:"for_symptom"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Source      string `json:"source,omitempty"`
}
