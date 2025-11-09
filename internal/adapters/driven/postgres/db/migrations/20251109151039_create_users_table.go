package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateUsersTable() *gormigrate.Migration {
	// Define the table structure as a temporary struct.
	// GORM will automatically map these fields to snake_case column names.
	// e.g., 'UUID' becomes 'uuid', 'CreatedAt' becomes 'created_at'.
	type User struct {
		// ID is the auto-incrementing primary key (bigint).
		ID uint `gorm:"primarykey"`

		// UUID is the universally unique identifier (e.g., for public-facing IDs).
		// We use postgres 'uuid' type and a unique index for fast lookups.
		UUID string `gorm:"type:uuid;uniqueIndex;not null"`

		// Name of the user.
		Name string `gorm:"type:varchar(255);not null"`

		// Email is used for login and must be unique.
		// A unique index is critical for fast 'WHERE email = ?' queries.
		Email string `gorm:"type:varchar(255);uniqueIndex;not null"`

		// Password stores the hashed password (e.g., Argon2id).
		Password string `gorm:"type:varchar(255);not null"`

		// CreatedAt and UpdatedAt are standard GORM timestamp fields.
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	return &gormigrate.Migration{
		ID: "20251109151039",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&User{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&User{})
		},
	}
}
