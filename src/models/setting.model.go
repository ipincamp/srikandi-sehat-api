package models

import "time"

type Setting struct {
	Key       string    `gorm:"primaryKey;type:varchar(100)"`
	Value     string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
