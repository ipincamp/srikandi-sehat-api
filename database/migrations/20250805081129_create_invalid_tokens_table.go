package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateInvalidTokensTable() *gormigrate.Migration {
	type InvalidToken struct {
		ID        uint      `gorm:"primarykey"`
		Token     string    `gorm:"type:text;uniqueIndex;not null"`
		ExpiresAt time.Time `gorm:"not null"`
	}

	return &gormigrate.Migration{
		ID: "20250805081129",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&InvalidToken{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&InvalidToken{})
		},
	}
}
