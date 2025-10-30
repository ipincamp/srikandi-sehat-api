package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddEmailVerification() *gormigrate.Migration {
	type User struct {
		EmailVerifiedAt       *time.Time `gorm:"type:datetime(3);null"`
		VerificationToken     string     `gorm:"type:varchar(255);null;index"`
		VerificationExpiresAt *time.Time `gorm:"type:datetime(3);null"`
	}

	return &gormigrate.Migration{
		ID: "20251030140355",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&User{})
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropColumn(&User{}, "email_verified_at"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn(&User{}, "verification_token"); err != nil {
				return err
			}
			return tx.Migrator().DropColumn(&User{}, "verification_expires_at")
		},
	}
}
