package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateUsersTable() *gormigrate.Migration {
	type User struct {
		ID        uint   `gorm:"primarykey"`
		Name      string `gorm:"type:varchar(100)"`
		Email     string `gorm:"type:varchar(255);uniqueIndex;not null"`
		Password  string `gorm:"type:varchar(255);not null"`
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	return &gormigrate.Migration{
		ID: "20250612000001",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&User{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&User{})
		},
	}
}
