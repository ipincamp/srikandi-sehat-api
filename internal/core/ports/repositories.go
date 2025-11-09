package ports

import (
	"context"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
)

// UserRepository is the "driven port" for user persistence.
// It defines the contract that our core services (e.g., UserService)
// will use to interact with the database, without knowing the
// specific database technology.
type UserRepository interface {
	// Save creates or updates a user in the data store.
	Save(ctx context.Context, user *domain.User) error

	// FindByID retrieves a user by their public UUID.
	FindByID(ctx context.Context, uuid string) (*domain.User, error)

	// FindByEmail retrieves a user by their email address.
	// This is crucial for login and registration checks.
	FindByEmail(ctx context.Context, email string) (*domain.User, error)

	// DEPRECATED
	// GetAllUserEmails retrieves all user emails from the data store.
	// This is used to populate in-memory caches on application startup.
	// GetAllUserEmails(ctx context.Context) ([]string, error)

	// TODO: Add other necessary methods like Delete, List, etc. as needed.
}

// NOTE: We would also define other repository interfaces here, e.g.:
// type TokenRepository interface {
//     SaveRefreshToken(ctx, userID, tokenID string) error
//     IsRefreshTokenValid(ctx, userID, tokenID string) bool
//     DeleteRefreshToken(ctx, userID, tokenID string) error
// }
