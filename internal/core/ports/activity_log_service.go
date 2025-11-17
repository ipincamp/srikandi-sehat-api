package ports

import (
	"context"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
)

// ActivityLogService defines the business logic for retrieving activity logs.
type ActivityLogService interface {
	GetMyLogs(
		ctx context.Context,
		userID string,
		page int,
		limit int,
		actionType *string,
	) (*domain.PaginatedActivityLogs, error)
}
