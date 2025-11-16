package resolvers

import (
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
