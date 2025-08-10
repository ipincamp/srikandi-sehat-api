package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primarykey"`
	UUID     string `gorm:"type:char(36);uniqueIndex;not null"`
	Name     string `gorm:"type:varchar(100)"`
	Email    string `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password string `gorm:"type:varchar(255);not null"`

	Roles       []*Role       `gorm:"many2many:user_roles;"`
	Permissions []*Permission `gorm:"many2many:user_permissions;"`
	Profile     Profile       `gorm:"foreignKey:UserID"`
	// MenstrualCycles []menstrual.MenstrualCycle `gorm:"foreignKey:UserID"`
	// SymptomLogs     []menstrual.SymptomLog     `gorm:"foreignKey:UserID"`

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
