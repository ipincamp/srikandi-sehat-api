package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primarykey"`
	UUID     string `gorm:"type:char(36);uniqueIndex;not null"`
	Name     string `gorm:"type:varchar(100)"`
	Email    string `gorm:"type:varchar(255);uniqueIndex;not null"`
	FcmToken string `gorm:"type:text"`

	EmailVerifiedAt       sql.NullTime   `gorm:"column:email_verified_at"`
	VerificationToken     sql.NullString `gorm:"column:verification_token"`
	VerificationExpiresAt sql.NullTime   `gorm:"column:verification_expires_at"`
	LastOTPSentAt         sql.NullTime   `gorm:"column:last_otp_sent_at"`

	Roles         []*Role            `gorm:"many2many:user_roles;"`
	Permissions   []*Permission      `gorm:"many2many:user_permissions;"`
	Profile       Profile            `gorm:"foreignKey:UserID"`
	AuthProviders []UserAuthProvider `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`

	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	if user.UUID == "" {
		user.UUID = uuid.New().String()
	}
	return
}
