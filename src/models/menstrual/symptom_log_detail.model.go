package menstrual

import "database/sql"

type SymptomLogDetail struct {
	ID uint `gorm:"primarykey"`

	SymptomLogID    uint       `gorm:"not null"`
	SymptomLog      SymptomLog `gorm:"foreignKey:SymptomLogID"`
	SymptomID       uint       `gorm:"not null"`
	Symptom         Symptom    `gorm:"foreignKey:SymptomID"`
	SymptomOptionID sql.NullInt64
	SymptomOption   SymptomOption `gorm:"foreignKey:SymptomOptionID"`
}
