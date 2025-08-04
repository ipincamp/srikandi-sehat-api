package models

import "time"

type User struct {
	ID       uint   `gorm:"primarykey"`
	Name     string `gorm:"type:varchar(100)"`
	Email    string `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password string `gorm:"type:varchar(255);not null"`

	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
