package models

import "time"

type Permission struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"type:varchar(100);uniqueIndex;not null"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
