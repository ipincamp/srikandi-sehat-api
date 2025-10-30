package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddLastOtpSentAt() *gormigrate.Migration {
	type User struct {
		LastOTPSentAt *time.Time `gorm:"column:last_otp_sent_at;type:datetime(3);null"`
	}

	return &gormigrate.Migration{
		ID: "20251030153032",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&User{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn(&User{}, "last_otp_sent_at")
		},
	}
}
