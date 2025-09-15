package menstrual

import "time"

type SymptomOption struct {
	ID    uint   `gorm:"primarykey"`
	Name  string `gorm:"type:varchar(100)"`
	Value string `gorm:"type:varchar(255)"`

	SymptomID uint    `gorm:"not null"`
	Symptom   Symptom `gorm:"foreignKey:SymptomID"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
