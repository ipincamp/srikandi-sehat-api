package region

import "time"

type Regency struct {
	ID   uint   `gorm:"primarykey"`
	Code string `gorm:"type:char(4);uniqueIndex"`
	Name string `gorm:"type:varchar(100)"`

	ProvinceID uint       `gorm:"not null"`
	Province   Province   `gorm:"foreignKey:ProvinceID"`
	Districts  []District `gorm:"foreignKey:RegencyID"`

	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
