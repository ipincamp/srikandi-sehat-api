package region

import "time"

type Classification struct {
	ID        uint   `gorm:"primarykey"`
	Name      string `gorm:"type:char(9);uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
