package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateRegionTables() *gormigrate.Migration {
	type Classification struct {
		ID        uint      `gorm:"primarykey"`
		Name      string    `gorm:"type:varchar(20);uniqueIndex"`
		CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
		UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	}

	type Province struct {
		ID        uint      `gorm:"primarykey"`
		Code      string    `gorm:"type:char(2);uniqueIndex"`
		Name      string    `gorm:"type:varchar(100)"`
		CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
		UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	}

	type Regency struct {
		ID         uint      `gorm:"primarykey"`
		Code       string    `gorm:"type:char(4);uniqueIndex"`
		Name       string    `gorm:"type:varchar(100)"`
		ProvinceID uint      `gorm:"not null"`
		CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
		UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	}

	type District struct {
		ID        uint      `gorm:"primarykey"`
		Code      string    `gorm:"type:char(7);uniqueIndex"`
		Name      string    `gorm:"type:varchar(100)"`
		RegencyID uint      `gorm:"not null"`
		CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
		UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	}

	type Village struct {
		ID               uint      `gorm:"primarykey"`
		Code             string    `gorm:"type:char(11);uniqueIndex"`
		Name             string    `gorm:"type:varchar(100)"`
		DistrictID       uint      `gorm:"not null"`
		ClassificationID uint      `gorm:"not null"`
		CreatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP"`
		UpdatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	}

	return &gormigrate.Migration{
		ID: "20250805083354",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				&Classification{},
				&Province{},
				&Regency{},
				&District{},
				&Village{},
			)
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(
				&Village{},
				&District{},
				&Regency{},
				&Province{},
				&Classification{},
			)
		},
	}
}
