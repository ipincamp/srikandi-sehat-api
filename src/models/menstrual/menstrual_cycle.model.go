package menstrual

import (
	"database/sql"
	"ipincamp/srikandi-sehat/src/models"
	"time"
)

type MenstrualCycle struct {
	ID             uint `gorm:"primarykey"`
	StartDate      time.Time
	EndDate        sql.NullTime
	PeriodLength   sql.NullInt16
	CycleLength    sql.NullInt16
	IsPeriodNormal sql.NullBool
	IsCycleNormal  sql.NullBool

	UserID      uint         `gorm:"not null"`
	User        models.User  `gorm:"foreignKey:UserID"`
	SymptomLogs []SymptomLog `gorm:"foreignKey:MenstrualCycleID"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
