package models

import "time"

type Notification struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Title     string    `gorm:"type:varchar(255);not null" json:"title"`
	Body      string    `gorm:"type:text" json:"body"`
	IsRead    bool      `gorm:"default:false" json:"is_read"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
