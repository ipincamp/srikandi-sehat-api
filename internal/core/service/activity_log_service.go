package service

import (
	"context"
	"errors"
	"math"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
	"github.com/rs/zerolog"
)

var _ ports.ActivityLogService = (*activityLogService)(nil)

type activityLogService struct {
	repo   ports.ActivityLogRepository
	logger zerolog.Logger
}

// NewActivityLogService is the constructor.
// Note: It uses the *non-transactional* repo, as this is for reads.
func NewActivityLogService(repo ports.ActivityLogRepository, logger zerolog.Logger) ports.ActivityLogService {
	return &activityLogService{
		repo:   repo,
		logger: logger,
	}
}

// GetMyLogs implements the logic for Requirement 1.14.
func (s *activityLogService) GetMyLogs(
	ctx context.Context,
	userID string,
	page int,
	limit int,
	actionType *string,
) (*domain.PaginatedActivityLogs, error) {
	log := s.logger.With().Str("method", "GetMyLogs").Str("user_id", userID).Logger()

	// 1. Set defaults for pagination (Req 1.14.3.2)
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	// 2. Calculate offset
	offset := (page - 1) * limit

	// 3. Query Total (Req 1.14.3.4)
	totalItems, err := s.repo.CountByUserID(ctx, userID, actionType)
	if err != nil {
		log.Error().Err(err).Msg("Failed to count user logs")
		return nil, errors.New("failed to retrieve logs")
	}

	var logs []*domain.ActivityLog
	var totalPages int = 0

	// 4. Query Logs if any exist (Req 1.14.3.3)
	if totalItems > 0 {
		logs, err = s.repo.FindByUserID(ctx, userID, actionType, offset, limit)
		if err != nil {
			log.Error().Err(err).Msg("Failed to find user logs")
			return nil, errors.New("failed to retrieve logs")
		}
		totalPages = int(math.Ceil(float64(totalItems) / float64(limit)))
	} else {
		// No records found, return empty slice
		logs = make([]*domain.ActivityLog, 0)
	}

	// 5. Assemble response (Req 1.14.3.5)
	response := &domain.PaginatedActivityLogs{
		Logs: logs,
		Pagination: domain.PaginationMetadata{
			TotalItems:  totalItems,
			TotalPages:  totalPages,
			CurrentPage: page,
			PerPage:     limit,
		},
	}

	return response, nil
}
