package ports

import (
	"context"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
)

// PersonalTokenRepository defines the persistence operations for personal tokens.
type PersonalTokenRepository interface {
	Save(ctx context.Context, token *domain.PersonalToken) error
	FindByID(ctx context.Context, jti string) (*domain.PersonalToken, error)
	Delete(ctx context.Context, jti string) error
	DeleteByUserID(ctx context.Context, userID string) error
}
