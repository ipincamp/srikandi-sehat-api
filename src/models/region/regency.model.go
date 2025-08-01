package region

import "time"

type Regency struct {
	ID         uint `gorm:"primarykey"`
	ProvinceID uint
	Code       string     `gorm:"type:char(4);uniqueIndex"`
	Name       string     `gorm:"type:varchar(100)"`
	Districts  []District `gorm:"foreignKey:RegencyID"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
