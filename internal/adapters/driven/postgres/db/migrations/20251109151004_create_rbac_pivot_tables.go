package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateRbacPivotTables() *gormigrate.Migration {
	type UserRole struct {
		UserID string `gorm:"type:uuid;primarykey"`
		RoleID string `gorm:"type:uuid;primarykey"`
	}

	type RolePermission struct {
		PermissionID string `gorm:"type:uuid;primarykey"`
		RoleID       string `gorm:"type:uuid;primarykey"`
	}

	type UserPermission struct {
		UserID       string `gorm:"type:uuid;primarykey"`
		PermissionID string `gorm:"type:uuid;primarykey"`
	}

	return &gormigrate.Migration{
		ID: "20251109151004",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&UserRole{}, &RolePermission{}, &UserPermission{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&UserPermission{}, &RolePermission{}, &UserRole{})
		},
	}
}
