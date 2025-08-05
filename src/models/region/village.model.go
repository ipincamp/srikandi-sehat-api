package region

import "time"

type Village struct {
	ID   uint   `gorm:"primarykey"`
	Code string `gorm:"type:char(11);uniqueIndex"`
	Name string `gorm:"type:varchar(100)"`

	DistrictID       uint           `gorm:"not null"`
	District         District       `gorm:"foreignKey:DistrictID"`
	ClassificationID uint           `gorm:"not null"`
	Classification   Classification `gorm:"foreignKey:ClassificationID"`

	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
