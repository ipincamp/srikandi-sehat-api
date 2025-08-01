package region

import "time"

type Village struct {
	ID               uint `gorm:"primarykey"`
	DistrictID       uint
	ClassificationID uint
	Code             string `gorm:"type:char(10);uniqueIndex"`
	Name             string `gorm:"type:varchar(100)"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
