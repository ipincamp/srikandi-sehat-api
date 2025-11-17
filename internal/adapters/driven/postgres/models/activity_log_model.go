package models

import (
	"time"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
)

// ActivityLogDBModel is the GORM model for the 'activity_logs' table.
// This matches the migration 20251109151006.
type ActivityLogDBModel struct {
	ID          int64     `gorm:"primarykey;autoIncrement"`
	ActorUserID *string   `gorm:"type:uuid;index"`
	Action      string    `gorm:"type:varchar(100);not null"`
	TargetTable *string   `gorm:"type:varchar(100);index,priority:1"`
	TargetID    *string   `gorm:"type:uuid;index,priority:2"`
	Changes     *string   `gorm:"type:jsonb"`
	IPAddress   *string   `gorm:"type:varchar(50)"`
	UserAgent   *string   `gorm:"type:text"`
	RequestID   *string   `gorm:"type:uuid"`
	Timestamp   time.Time `gorm:"type:timestamptz;not null;default:NOW()"`
}

// TableName specifies the table name for GORM.
func (ActivityLogDBModel) TableName() string {
	return "activity_logs"
}

// ToDomain converts the DB model to the domain model.
func (m *ActivityLogDBModel) ToDomain() *domain.ActivityLog {
	return &domain.ActivityLog{
		ID:          m.ID,
		ActorUserID: m.ActorUserID,
		Action:      m.Action,
		TargetTable: m.TargetTable,
		TargetID:    m.TargetID,
		Changes:     m.Changes,
		IPAddress:   m.IPAddress,
		UserAgent:   m.UserAgent,
		RequestID:   m.RequestID,
		Timestamp:   m.Timestamp,
	}
}

// FromDomain converts the domain model to the DB model.
func ActivityLogFromDomain(d *domain.ActivityLog) *ActivityLogDBModel {
	return &ActivityLogDBModel{
		ID:          d.ID,
		ActorUserID: d.ActorUserID,
		Action:      d.Action,
		TargetTable: d.TargetTable,
		TargetID:    d.TargetID,
		Changes:     d.Changes,
		IPAddress:   d.IPAddress,
		UserAgent:   d.UserAgent,
		RequestID:   d.RequestID,
		Timestamp:   d.Timestamp,
	}
}
