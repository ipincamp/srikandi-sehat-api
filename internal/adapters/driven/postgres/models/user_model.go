package models

import (
	"time"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
	"gorm.io/gorm"
)

// UserDBModel merepresentasikan skema tabel 'users' di database.
type UserDBModel struct {
	ID              string `gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
	Name            string `gorm:"type:varchar(100);not null"`
	Email           string `gorm:"type:varchar(255);not null;index:idx_users_email_active_unique,where:deleted_at IS NULL"`
	EmailVerifiedAt *time.Time
	PasswordHash    string `gorm:"type:varchar(255);not null"`
	DisabledAt      *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

// TableName menentukan nama tabel secara eksplisit
func (UserDBModel) TableName() string {
	return "users"
}

// ToDomain mengkonversi model DB (infrastruktur) ke model Domain (bisnis).
func (m *UserDBModel) ToDomain() *domain.User {
	return &domain.User{
		ID:              m.ID,
		Name:            m.Name,
		Email:           m.Email,
		EmailVerifiedAt: m.EmailVerifiedAt,
		PasswordHash:    m.PasswordHash,
		DisabledAt:      m.DisabledAt,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		DeletedAt:       &m.DeletedAt.Time,
	}
}

// FromDomain mengkonversi model Domain (bisnis) ke model DB (infrastruktur).
func FromDomain(d *domain.User) *UserDBModel {
	var deletedAt gorm.DeletedAt
	if d.DeletedAt != nil {
		deletedAt.Time = *d.DeletedAt
		deletedAt.Valid = true
	}

	return &UserDBModel{
		ID:              d.ID,
		Name:            d.Name,
		Email:           d.Email,
		EmailVerifiedAt: d.EmailVerifiedAt,
		PasswordHash:    d.PasswordHash,
		DisabledAt:      d.DisabledAt,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
		DeletedAt:       deletedAt,
	}
}
