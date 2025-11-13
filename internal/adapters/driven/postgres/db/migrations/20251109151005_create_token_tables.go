package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateTokenTables() *gormigrate.Migration {
	type PersonalToken struct {
		ID        string    `gorm:"type:uuid;primarykey"`
		UserID    string    `gorm:"type:uuid;not null;index"`
		ExpiresAt time.Time `gorm:"type:timestamptz;not null"`
		CreatedAt time.Time `gorm:"type:timestamptz;not null;default:NOW()"`
	}

	type UserToken struct {
		UserID    string    `gorm:"type:uuid;primarykey"`
		Purpose   string    `gorm:"type:varchar(50);primarykey"`
		TokenHash string    `gorm:"type:varchar(255);not null;index"`
		ExpiresAt time.Time `gorm:"type:timestamptz;not null"`
		CreatedAt time.Time `gorm:"type:timestamptz;not null;default:NOW()"`
	}

	return &gormigrate.Migration{
		ID: "20251109151005",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&PersonalToken{}, &UserToken{})
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropTable(&PersonalToken{}, &UserToken{}); err != nil {
				return err
			}
			return nil
		},
	}
}
