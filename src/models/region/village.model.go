package region

import "time"

type Village struct {
	ID   uint   `gorm:"primarykey"`
	Code string `gorm:"type:char(10);uniqueIndex"`
	Name string `gorm:"type:varchar(100)"`

	DistrictID       uint
	District         District `gorm:"foreignKey:DistrictID"`
	ClassificationID uint
	Classification   Classification `gorm:"foreignKey:ClassificationID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
