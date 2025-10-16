package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddFcmTokenToUsersTable() *gormigrate.Migration {
	type User struct {
		FcmToken string `gorm:"type:text"`
	}

	return &gormigrate.Migration{
		ID: "20251016183214",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&User{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn(&User{}, "fcm_token")
		},
	}
}
