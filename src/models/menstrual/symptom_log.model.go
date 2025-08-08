package menstrual

import (
	"database/sql"
	"time"
)

type SymptomLog struct {
	ID       uint      `gorm:"primarykey"`
	LoggedAt time.Time `gorm:"column:logged_at"`
	Note     string    `gorm:"type:text"`

	UserID           uint               `gorm:"not null;index"`
	Details          []SymptomLogDetail `gorm:"foreignKey:SymptomLogID"`
	MenstrualCycleID sql.NullInt64

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
