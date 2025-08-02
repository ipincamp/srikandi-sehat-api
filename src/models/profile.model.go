package models

import (
	"ipincamp/srikandi-sehat/src/constants"
	"time"
)

type Profile struct {
	ID                  uint       `gorm:"primarykey"`
	UserID              uint       `gorm:"uniqueIndex;not null"`
	PhotoURL            string     `gorm:"type:varchar(255)"`
	PhoneNumber         string     `gorm:"type:varchar(20);uniqueIndex"`
	DateOfBirth         *time.Time `gorm:"type:date"`
	HeightCM            uint       `gorm:"type:smallint"`
	WeightKG            float32    `gorm:"type:decimal(5,2)"`
	AddressStreet       string     `gorm:"type:text"`
	VillageID           *uint
	LastEducation       constants.EducationLevel `gorm:"type:enum('Tidak Sekolah', 'SD', 'SMP', 'SMA', 'Diploma', 'S1', 'S2', 'S3')"`
	ParentLastEducation constants.EducationLevel `gorm:"type:enum('Tidak Sekolah', 'SD', 'SMP', 'SMA', 'Diploma', 'S1', 'S2', 'S3')"`
	ParentLastJob       string                   `gorm:"type:varchar(100)"`
	InternetAccess      constants.InternetAccess `gorm:"type:enum('WiFi', 'Seluler')"`
	MenarcheAge         uint                     `gorm:"type:tinyint"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
