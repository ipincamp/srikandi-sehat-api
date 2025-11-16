package ports

import (
	"context"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
)

// UserTokenRepository defines the persistence operations for one-time user tokens.
// This is used for email verification, password resets, etc.
type UserTokenRepository interface {
	// Save creates or updates a token (Upsert).
	// This matches the logic from the requirements.
	Save(ctx context.Context, token *domain.UserToken) error

	// FindByUserIDAndPurpose finds a specific token.
	FindByUserIDAndPurpose(ctx context.Context, userID string, purpose string) (*domain.UserToken, error)

	// DeleteByUserIDAndPurpose removes a specific token.
	DeleteByUserIDAndPurpose(ctx context.Context, userID string, purpose string) error
}
