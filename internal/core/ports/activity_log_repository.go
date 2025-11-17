package ports

import (
	"context"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
)

// ActivityLogRepository defines the persistence operations for activity logs.
type ActivityLogRepository interface {
	// Save creates a new log entry.
	Save(ctx context.Context, log *domain.ActivityLog) error

	// FindByUserID retrieves a paginated list of logs for a user.
	FindByUserID(
		ctx context.Context,
		userID string,
		actionType *string,
		offset int,
		limit int,
	) ([]*domain.ActivityLog, error)

	// CountByUserID counts the total logs for a user, for pagination.
	CountByUserID(ctx context.Context, userID string, actionType *string) (int64, error)
}
