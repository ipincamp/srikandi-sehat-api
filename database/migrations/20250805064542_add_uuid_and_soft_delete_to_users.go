package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddUuidAndSoftDeleteToUsers() *gormigrate.Migration {
	type User struct {
		ID uint `gorm:"primarykey"`
	}

	return &gormigrate.Migration{
		ID: "20250805064542",

		Migrate: func(tx *gorm.DB) error {
			tableName := tx.NamingStrategy.TableName("User")

			if err := tx.Exec("ALTER TABLE " + tableName + " ADD COLUMN uuid CHAR(36) NOT NULL AFTER id").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE " + tableName + " ADD UNIQUE INDEX idx_users_uuid (uuid)").Error; err != nil {
				return err
			}

			if err := tx.Exec("ALTER TABLE " + tableName + " ADD COLUMN deleted_at DATETIME(3) NULL AFTER updated_at").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE " + tableName + " ADD INDEX idx_users_deleted_at (deleted_at)").Error; err != nil {
				return err
			}
			return nil
		},

		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropColumn(&User{}, "uuid"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn(&User{}, "deleted_at"); err != nil {
				return err
			}
			return nil
		},
	}
}
