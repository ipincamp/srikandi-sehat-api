package region

import "time"

type District struct {
	ID        uint `gorm:"primarykey"`
	RegencyID uint
	Code      string    `gorm:"type:char(7);uniqueIndex"`
	Name      string    `gorm:"type:varchar(100)"`
	Villages  []Village `gorm:"foreignKey:DistrictID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
