package dto

import "time"

// FullExportRecord defines the complete flattened structure for the combined user and cycle data CSV export.
type FullExportRecord struct {
	// User Data
	UserUUID            string    `json:"user_uuid"`
	UserName            string    `json:"user_name"`
	UserEmail           string    `json:"user_email"`
	UserRegisteredAt    time.Time `json:"user_registered_at"`
	Age                 int       `json:"age"`
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

	// Cycle Data
	CycleNumber    int64  `json:"cycle_number"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	PeriodLength   int16  `json:"period_length"`
	PeriodCategory string `json:"period_category"`
	CycleLength    int16  `json:"cycle_length"`
	CycleCategory  string `json:"cycle_category"`
	Symptoms       string `json:"symptoms"`
}
