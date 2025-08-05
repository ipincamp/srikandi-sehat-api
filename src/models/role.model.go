package models

import "time"

type Role struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"type:varchar(100);uniqueIndex;not null"`

	Permissions []*Permission `gorm:"many2many:role_permissions;"`

	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
