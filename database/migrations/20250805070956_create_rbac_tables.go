package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateRbacTables() *gormigrate.Migration {
	type Permission struct {
		ID        uint      `gorm:"primarykey"`
		Name      string    `gorm:"type:varchar(100);uniqueIndex;not null"`
		CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
		UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	}

	type Role struct {
		ID          uint          `gorm:"primarykey"`
		Name        string        `gorm:"type:varchar(100);uniqueIndex;not null"`
		Permissions []*Permission `gorm:"many2many:role_permissions;"`
		CreatedAt   time.Time     `gorm:"default:CURRENT_TIMESTAMP"`
		UpdatedAt   time.Time     `gorm:"default:CURRENT_TIMESTAMP"`
	}

	type User struct {
		ID          uint
		Roles       []*Role       `gorm:"many2many:user_roles;"`
		Permissions []*Permission `gorm:"many2many:user_permissions;"`
	}

	return &gormigrate.Migration{
		ID: "20250805070956",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				&Permission{},
				&Role{},
				&User{},
			)
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(
				"user_roles",
				"user_permissions",
				"role_permissions",
				&Role{},
				&Permission{},
			)
		},
	}
}
