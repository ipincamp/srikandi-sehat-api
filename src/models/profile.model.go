package models

import (
	"ipincamp/srikandi-sehat/src/constants"
	"ipincamp/srikandi-sehat/src/models/region"
	"time"
)

type Profile struct {
	ID                  uint                     `gorm:"primarykey"`
	PhotoURL            string                   `gorm:"type:varchar(255)"`
	PhoneNumber         string                   `gorm:"type:varchar(20)"`
	DateOfBirth         *time.Time               `gorm:"type:date"`
	HeightCM            uint                     `gorm:"type:smallint"`
	WeightKG            float32                  `gorm:"type:decimal(5,2)"`
	LastEducation       constants.EducationLevel `gorm:"type:enum('Tidak Sekolah', 'SD', 'SMP', 'SMA', 'Diploma', 'S1', 'S2', 'S3')"`
	ParentLastEducation constants.EducationLevel `gorm:"type:enum('Tidak Sekolah', 'SD', 'SMP', 'SMA', 'Diploma', 'S1', 'S2', 'S3')"`
	ParentLastJob       string                   `gorm:"type:varchar(100)"`
	InternetAccess      constants.InternetAccess `gorm:"type:enum('WiFi', 'Seluler')"`
	MenarcheAge         uint                     `gorm:"type:tinyint"`

	UserID    uint `gorm:"uniqueIndex;not null"`
	VillageID *uint
	Village   region.Village `gorm:"foreignKey:VillageID"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
