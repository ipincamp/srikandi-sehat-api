package region

import "time"

type Province struct {
	ID   uint   `gorm:"primarykey"`
	Code string `gorm:"type:char(2);uniqueIndex"`
	Name string `gorm:"type:varchar(100)"`

	Regencies []Regency `gorm:"foreignKey:ProvinceID"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
