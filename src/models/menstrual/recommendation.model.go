package menstrual

import "time"

type Recommendation struct {
	ID          uint   `gorm:"primarykey"`
	Title       string `gorm:"type:varchar(255)"`
	Description string `gorm:"type:text"`
	Source      string `gorm:"type:varchar(255)"`

	SymptomID uint    `gorm:"not null"`
	Symptom   Symptom `gorm:"foreignKey:SymptomID"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
