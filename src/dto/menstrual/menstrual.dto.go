package dto

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
