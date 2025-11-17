package postgres

import (
	"context"

	"github.com/ipincamp/srikandi-sehat/internal/adapters/driven/postgres/models"
	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

var _ ports.ActivityLogRepository = (*activityLogRepository)(nil)

type activityLogRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

// NewActivityLogRepository is the constructor.
func NewActivityLogRepository(db *gorm.DB, logger zerolog.Logger) ports.ActivityLogRepository {
	return &activityLogRepository{
		db:     db,
		logger: logger,
	}
}

// Save creates a new log entry.
func (r *activityLogRepository) Save(ctx context.Context, log *domain.ActivityLog) error {
	model := models.ActivityLogFromDomain(log)
	// We don't want GORM to touch the Timestamp, it's set by default in DB
	return r.db.WithContext(ctx).Omit("Timestamp").Create(model).Error
}

// buildUserLogQuery is a helper to build the base query.
func (r *activityLogRepository) buildUserLogQuery(ctx context.Context, userID string, actionType *string) *gorm.DB {
	query := r.db.WithContext(ctx).Model(&models.ActivityLogDBModel{}).Where("actor_user_id = ?", userID)

	if actionType != nil && *actionType != "" {
		query = query.Where("action = ?", *actionType)
	}
	return query
}

// FindByUserID retrieves a paginated list of logs (Req 1.14.3.3).
func (r *activityLogRepository) FindByUserID(
	ctx context.Context,
	userID string,
	actionType *string,
	offset int,
	limit int,
) ([]*domain.ActivityLog, error) {
	var models []models.ActivityLogDBModel

	query := r.buildUserLogQuery(ctx, userID, actionType)

	err := query.Order("timestamp DESC").Offset(offset).Limit(limit).Find(&models).Error
	if err != nil {
		return nil, err
	}

	// Convert models to domain
	logs := make([]*domain.ActivityLog, len(models))
	for i, m := range models {
		logs[i] = m.ToDomain()
	}
	return logs, nil
}

// CountByUserID counts the total logs for a user (Req 1.14.3.4).
func (r *activityLogRepository) CountByUserID(ctx context.Context, userID string, actionType *string) (int64, error) {
	var count int64
	query := r.buildUserLogQuery(ctx, userID, actionType)
	err := query.Count(&count).Error
	return count, err
}
