package migrations

import (
	"ipincamp/srikandi-sehat/src/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateUsersTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20250612000001",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.User{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&models.User{})
		},
	}
}
