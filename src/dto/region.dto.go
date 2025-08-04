package dto

// Request Query
type RegencyQuery struct {
	ProvinceCode string `query:"province_code" validate:"required,len=2"`
}

type DistrictQuery struct {
	RegencyCode string `query:"regency_code" validate:"required,len=4"`
}

type VillageQuery struct {
	DistrictCode string `query:"district_code" validate:"required,len=7"`
}

// Response Body
type RegionResponse struct {
	Code           string `json:"code"`
	Name           string `json:"name"`
	Classification string `json:"type,omitempty"`
}
