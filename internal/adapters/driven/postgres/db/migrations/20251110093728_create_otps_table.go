package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateOtpsTable() *gormigrate.Migration {
	type OTP struct {
		ID        uint      `gorm:"primarykey"`
		UserID    *uint     `gorm:"index"` // Nullable
		Email     string    `gorm:"type:varchar(255);index"`
		Code      string    `gorm:"type:varchar(10);index;not null"`
		Type      string    `gorm:"type:varchar(50);not null"`
		ExpiresAt time.Time `gorm:"not null"`
		CreatedAt time.Time
		// Tambahkan foreign key constraint jika Anda mau
		// User      User      `gorm:"foreignKey:UserID;references:ID;OnDelete:CASCADE"`
	}

	return &gormigrate.Migration{
		ID: "20251110093728",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&OTP{})

		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&OTP{})

		},
	}
}
