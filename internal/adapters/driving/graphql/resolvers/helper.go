package resolvers

import (
	"fmt"
	"time"

	"github.com/ipincamp/srikandi-sehat/internal/adapters/driving/graphql/models"
	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
)

func mapDomainUserToGqlUser(user *domain.User) *models.User {
	return &models.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}

func mapDomainToGqlLog(log *domain.ActivityLog) *models.ActivityLog {
	return &models.ActivityLog{
		ID:          fmt.Sprintf("%d", log.ID),
		Action:      log.Action,
		TargetTable: log.TargetTable,
		TargetID:    log.TargetID,
		Changes:     log.Changes,
		IPAddress:   log.IPAddress,
		UserAgent:   log.UserAgent,
		Timestamp:   log.Timestamp.Format(time.RFC3339),
	}
}

func mapDomainToGqlPaginatedLogs(paginated *domain.PaginatedActivityLogs) *models.PaginatedActivityLogs {
	gqlLogs := make([]*models.ActivityLog, len(paginated.Logs))
	for i, log := range paginated.Logs {
		gqlLogs[i] = mapDomainToGqlLog(log)
	}

	return &models.PaginatedActivityLogs{
		Data: gqlLogs,
		Meta: &models.ActivityLogPagination{
			TotalItems:  int(paginated.Pagination.TotalItems),
			TotalPages:  paginated.Pagination.TotalPages,
			CurrentPage: paginated.Pagination.CurrentPage,
			PerPage:     paginated.Pagination.PerPage,
		},
	}
}
