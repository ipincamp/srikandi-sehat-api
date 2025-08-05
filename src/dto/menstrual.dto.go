package dto

import (
	"database/sql"
	"ipincamp/srikandi-sehat/src/constants"
	"time"
)

// Request Query
type SymptomLogQuery struct {
	StartDate string `query:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate   string `query:"finish_date" validate:"required,datetime=2006-01-02"`
}

// Request Body
type CycleRequest struct {
	StartDate string `json:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate   string `json:"finish_date" validate:"omitempty,datetime=2006-01-02"`
}

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

type CycleResponse struct {
	ID             uint          `json:"id"`
	StartDate      time.Time     `json:"start_date"`
	EndDate        sql.NullTime  `json:"finish_date"`
	PeriodLength   sql.NullInt16 `json:"period_length"`
	CycleLength    sql.NullInt16 `json:"cycle_length"`
	IsPeriodNormal sql.NullBool  `json:"is_period_normal"`
	IsCycleNormal  sql.NullBool  `json:"is_cycle_normal"`
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
