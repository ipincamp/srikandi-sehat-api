package region

import "time"

type Classification struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"type:varchar(20);uniqueIndex"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
