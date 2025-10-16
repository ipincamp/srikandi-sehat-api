package models

import "time"

type Notification struct {
	ID        uint      `gorm:"primarykey"`
	Title     string    `gorm:"type:varchar(255);not null"`
	Body      string    `gorm:"type:text"`
	IsRead    bool      `gorm:"default:false"`
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
