package models

import (
	"ipincamp/srikandi-sehat/src/constants"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uint                 `gorm:"primarykey" json:"-"`
	UUID     string               `gorm:"type:char(36);uniqueIndex" json:"id"`
	Name     string               `gorm:"type:varchar(255);not null" json:"name" validate:"required,min=3"`
	Email    string               `gorm:"type:varchar(255);uniqueIndex;not null" json:"email" validate:"required,email"`
	Password string               `gorm:"type:varchar(255);not null" json:"-"`
	Status   constants.UserStatus `gorm:"type:enum('processing','active','suspended');default:'processing'" json:"-"`

	Roles       []*Role       `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	Permissions []*Permission `gorm:"many2many:user_permissions;" json:"permissions,omitempty"`
	Profile     Profile       `gorm:"foreignKey:UserID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.UUID = uuid.New().String()
	return
}
