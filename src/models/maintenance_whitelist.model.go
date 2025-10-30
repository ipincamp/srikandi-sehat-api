package models

import "time"

type MaintenanceWhitelist struct {
	UserID    uint      `gorm:"primaryKey"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
