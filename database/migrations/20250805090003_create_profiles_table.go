package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateProfilesTable() *gormigrate.Migration {
	type Profile struct {
		ID                  uint       `gorm:"primarykey"`
		PhotoURL            string     `gorm:"type:varchar(255)"`
		PhoneNumber         string     `gorm:"type:varchar(20)"`
		DateOfBirth         *time.Time `gorm:"type:date"`
		HeightCM            uint       `gorm:"type:smallint"`
		WeightKG            float32    `gorm:"type:decimal(5,2)"`
		LastEducation       string     `gorm:"type:enum('Tidak Sekolah', 'SD', 'SMP', 'SMA', 'Diploma', 'S1', 'S2', 'S3')"`
		ParentLastEducation string     `gorm:"type:enum('Tidak Sekolah', 'SD', 'SMP', 'SMA', 'Diploma', 'S1', 'S2', 'S3')"`
		ParentLastJob       string     `gorm:"type:varchar(100)"`
		InternetAccess      string     `gorm:"type:enum('WiFi', 'Seluler')"`
		MenarcheAge         uint       `gorm:"type:tinyint"`
		UserID              uint       `gorm:"uniqueIndex;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		VillageID           *uint      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		CreatedAt           time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
		UpdatedAt           time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	}

	return &gormigrate.Migration{
		ID: "20250805090003",

		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&Profile{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&Profile{})
		},
	}
}
