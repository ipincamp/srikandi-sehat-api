package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreatePermissionsTable() *gormigrate.Migration {
	type Permission struct {
		ID          string `gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
		Name        string `gorm:"type:varchar(255);uniqueIndex;not null"`
		Description string `gorm:"type:text"`
	}

	return &gormigrate.Migration{
		ID: "20251109151002",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&Permission{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&Permission{})
		},
	}
}
