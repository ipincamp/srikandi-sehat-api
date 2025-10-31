package migrations

import (
	"database/sql"
	"ipincamp/srikandi-sehat/src/models"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// --- Helper Structs & Methods ---

type UserAuthProviderMigration struct {
	ID         uint
	UserID     uint           `gorm:"not null;uniqueIndex:idx_user_provider"`
	Provider   string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_user_provider"`
	ProviderID string         `gorm:"type:varchar(255);null;uniqueIndex:idx_provider_id"`
	Password   sql.NullString `gorm:"type:varchar(255);null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (UserAuthProviderMigration) TableName() string {
	return "user_auth_providers"
}

type UserOldRollback struct {
	ID       uint
	Password sql.NullString `gorm:"type:varchar(255);null"`
}

func (UserOldRollback) TableName() string {
	return "users"
}

type UserOldFinalRollback struct {
	Password string `gorm:"type:varchar(255);not null"`
}

func (UserOldFinalRollback) TableName() string {
	return "users"
}

func RefactorAuthTables() *gormigrate.Migration {
	userAuthProviderTable := "user_auth_providers"

	return &gormigrate.Migration{
		ID: "20251031070450",
		Migrate: func(tx *gorm.DB) error {
			// 1. Buat tabel baru menggunakan struct top-level
			if err := tx.AutoMigrate(&UserAuthProviderMigration{}); err != nil {
				return err
			}

			// 2. Migrasi data password dari 'users' ke 'user_auth_providers'
			execSQL := `
				INSERT INTO user_auth_providers (user_id, provider, password, created_at, updated_at)
				SELECT id, 'local', password, created_at, updated_at
				FROM users
				WHERE password IS NOT NULL AND password != ''
			`
			if err := tx.Exec(execSQL).Error; err != nil {
				return err
			}

			// 3. Hapus kolom password dari tabel 'users'
			if tx.Migrator().HasColumn(&models.User{}, "password") {
				if err := tx.Migrator().DropColumn(&models.User{}, "password"); err != nil {
					return err
				}
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// 1. Tambahkan kembali kolom password ke 'users' (nullable)
			if !tx.Migrator().HasColumn(&UserOldRollback{}, "password") {
				if err := tx.Migrator().AddColumn(&UserOldRollback{}, "password"); err != nil {
					return err
				}
			}

			// 2. Kembalikan data password dari 'user_auth_providers' ke 'users'
			execSQL := `
				UPDATE users u
				JOIN user_auth_providers a ON u.id = a.user_id
				SET u.password = a.password
				WHERE a.provider = 'local' AND a.password IS NOT NULL
			`
			if err := tx.Exec(execSQL).Error; err != nil {
				return err
			}

			// 3. Hapus tabel 'user_auth_providers'
			if err := tx.Migrator().DropTable(userAuthProviderTable); err != nil {
				return err
			}

			// 4. Ubah kolom password kembali ke NOT NULL
			return tx.Migrator().AlterColumn(&UserOldFinalRollback{}, "password")
		},
	}
}
