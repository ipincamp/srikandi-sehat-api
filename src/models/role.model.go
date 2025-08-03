package models

import "time"

type Role struct {
	ID   uint   `gorm:"primarykey" json:"id"`
	Name string `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`

	Permissions []*Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
