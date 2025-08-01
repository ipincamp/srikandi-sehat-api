package models

import "time"

type InvalidToken struct {
	ID        uint      `gorm:"primarykey"`
	Token     string    `gorm:"type:text;uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
}
