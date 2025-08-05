package migrations

import (
	"ipincamp/srikandi-sehat/src/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddUuidAndSoftDeleteToUsers() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20250805064542",

		Migrate: func(tx *gorm.DB) error {
			if err := tx.Migrator().AddColumn(&models.User{}, "UUID"); err != nil {
				return err
			}
			if err := tx.Migrator().AddColumn(&models.User{}, "DeletedAt"); err != nil {
				return err
			}
			return nil
		},

		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropColumn(&models.User{}, "UUID"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn(&models.User{}, "DeletedAt"); err != nil {
				return err
			}
			return nil
		},
	}
}
