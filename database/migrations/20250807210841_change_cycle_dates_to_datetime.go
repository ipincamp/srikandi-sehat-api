package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func ChangeCycleDatesToDatetime() *gormigrate.Migration {
	type MenstrualCycle struct{}

	return &gormigrate.Migration{
		ID: "20250807210841",

		Migrate: func(tx *gorm.DB) error {
			tableName := tx.NamingStrategy.TableName("MenstrualCycle")
			if err := tx.Exec("ALTER TABLE " + tableName + " MODIFY COLUMN start_date DATETIME(3) NOT NULL").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE " + tableName + " MODIFY COLUMN end_date DATETIME(3) NULL").Error; err != nil {
				return err
			}
			return nil
		},

		Rollback: func(tx *gorm.DB) error {
			tableName := tx.NamingStrategy.TableName("MenstrualCycle")
			if err := tx.Exec("ALTER TABLE " + tableName + " MODIFY COLUMN start_date DATE NOT NULL").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE " + tableName + " MODIFY COLUMN end_date DATE NULL").Error; err != nil {
				return err
			}
			return nil
		},
	}
}
