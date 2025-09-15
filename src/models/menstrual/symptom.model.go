package menstrual

import (
	"ipincamp/srikandi-sehat/src/constants"
	"time"
)

type Symptom struct {
	ID       uint                  `gorm:"primarykey"`
	Name     string                `gorm:"type:varchar(100);uniqueIndex"`
	Category string                `gorm:"type:varchar(100)"`
	Type     constants.SymptomType `gorm:"type:enum('BASIC','OPTIONS');default:'BASIC'"`

	Options []SymptomOption `gorm:"foreignKey:SymptomID"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
