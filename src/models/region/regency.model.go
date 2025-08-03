package region

import "time"

type Regency struct {
	ID   uint   `gorm:"primarykey"`
	Code string `gorm:"type:char(4);uniqueIndex"`
	Name string `gorm:"type:varchar(255)"`

	ProvinceID uint
	Province   Province   `gorm:"foreignKey:ProvinceID"`
	Districts  []District `gorm:"foreignKey:RegencyID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
