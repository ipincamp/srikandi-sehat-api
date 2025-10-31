package models

import (
	"database/sql"
	"time"
)

type UserAuthProvider struct {
	ID         uint   `gorm:"primarykey"`
	UserID     uint   `gorm:"not null;uniqueIndex:idx_user_provider"`
	User       User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Provider   string `gorm:"type:varchar(20);not null;uniqueIndex:idx_user_provider"` // "local", "google"
	ProviderID string `gorm:"type:varchar(255);null;uniqueIndex:idx_provider_id"`      // ID unik dari Google

	Password sql.NullString `gorm:"type:varchar(255);null"` // (hanya untuk provider 'local')

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
