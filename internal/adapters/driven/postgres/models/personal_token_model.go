package models

import (
	"time"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
)

// PersonalTokenDBModel is the GORM model for the 'personal_tokens' table.
type PersonalTokenDBModel struct {
	ID        string    `gorm:"type:uuid;primarykey"`
	UserID    string    `gorm:"type:uuid;not null;index"`
	ExpiresAt time.Time `gorm:"type:timestamptz;not null"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:NOW()"`
}

// TableName specifies the table name for GORM.
func (PersonalTokenDBModel) TableName() string {
	return "personal_tokens"
}

// ToDomain converts the DB model to the domain model.
func (m *PersonalTokenDBModel) ToDomain() *domain.PersonalToken {
	return &domain.PersonalToken{
		ID:        m.ID,
		UserID:    m.UserID,
		ExpiresAt: m.ExpiresAt,
		CreatedAt: m.CreatedAt,
	}
}

// FromDomain converts the domain model to the DB model.
func PersonalTokenFromDomain(d *domain.PersonalToken) *PersonalTokenDBModel {
	return &PersonalTokenDBModel{
		ID:        d.ID,
		UserID:    d.UserID,
		ExpiresAt: d.ExpiresAt,
		CreatedAt: d.CreatedAt,
	}
}
