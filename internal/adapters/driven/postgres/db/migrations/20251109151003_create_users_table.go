package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateUsersTable() *gormigrate.Migration {
	type User struct {
		ID              string         `gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
		Name            string         `gorm:"type:varchar(100);not null"`
		Email           string         `gorm:"type:varchar(255);not null"`
		EmailVerifiedAt *time.Time     `gorm:"type:timestamptz"`
		PasswordHash    string         `gorm:"type:varchar(255);not null"`
		DisabledAt      *time.Time     `gorm:"type:timestamptz"`
		CreatedAt       time.Time      `gorm:"type:timestamptz;not null;default:NOW()"`
		UpdatedAt       time.Time      `gorm:"type:timestamptz;not null;default:NOW()"`
		DeletedAt       gorm.DeletedAt `gorm:"type:timestamptz"`
	}

	return &gormigrate.Migration{
		ID: "20251109151003",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(&User{}); err != nil {
				return err
			}
			return tx.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS "idx_users_email_active_unique" ON "users"("email") WHERE ("deleted_at" IS NULL);`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			tx.Exec(`DROP INDEX IF EXISTS "idx_users_email_active_unique"`)
			return tx.Migrator().DropTable(&User{})
		},
	}
}
