package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateRolesTable() *gormigrate.Migration {
	type Role struct {
		ID          string `gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
		Name        string `gorm:"type:varchar(100);uniqueIndex;not null"`
		Description string `gorm:"type:text"`
	}

	return &gormigrate.Migration{
		ID: "20251109151001",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`).Error; err != nil {
				return err
			}

			return tx.AutoMigrate(&Role{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&Role{})
		},
	}
}
