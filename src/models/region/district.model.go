package region

import "time"

type District struct {
	ID   uint   `gorm:"primarykey"`
	Code string `gorm:"type:char(7);uniqueIndex"`
	Name string `gorm:"type:varchar(100)"`

	RegencyID uint
	Regency   Regency   `gorm:"foreignKey:RegencyID"`
	Villages  []Village `gorm:"foreignKey:DistrictID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
