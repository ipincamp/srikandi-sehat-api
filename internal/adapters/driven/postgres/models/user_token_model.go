package models

import (
	"time"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
)

// UserTokenDBModel is the GORM model for the 'user_tokens' table.
// This structure is based on the migration file.
type UserTokenDBModel struct {
	UserID    string    `gorm:"type:uuid;primarykey"`
	Purpose   string    `gorm:"type:varchar(50);primarykey"`
	TokenHash string    `gorm:"type:varchar(255);not null;index"`
	ExpiresAt time.Time `gorm:"type:timestamptz;not null"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:NOW()"`
}

// TableName specifies the table name for GORM.
func (UserTokenDBModel) TableName() string {
	return "user_tokens"
}

// ToDomain converts the DB model to the domain model.
func (m *UserTokenDBModel) ToDomain() *domain.UserToken {
	return &domain.UserToken{
		UserID:    m.UserID,
		Purpose:   m.Purpose,
		TokenHash: m.TokenHash,
		ExpiresAt: m.ExpiresAt,
		CreatedAt: m.CreatedAt,
	}
}

// FromDomain converts the domain model to the DB model.
func UserTokenFromDomain(d *domain.UserToken) *UserTokenDBModel {
	return &UserTokenDBModel{
		UserID:    d.UserID,
		Purpose:   d.Purpose,
		TokenHash: d.TokenHash,
		ExpiresAt: d.ExpiresAt,
		CreatedAt: d.CreatedAt, // GORM handles default on create
	}
}
