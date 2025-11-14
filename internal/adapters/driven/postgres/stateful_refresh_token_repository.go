package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
)

// Compile-time interface implementation check.
var _ ports.StatefulRefreshTokenRepository = (*statefulRefreshTokenRepository)(nil)

// statefulRefreshTokenRepository implements ports.StatefulRefreshTokenRepository interface using PostgreSQL.
type statefulRefreshTokenRepository struct {
	db     dbExecutor
	logger zerolog.Logger
}

// NewStatefulRefreshTokenRepository creates a new stateful refresh token repository instance.
func NewStatefulRefreshTokenRepository(db dbExecutor, logger zerolog.Logger) ports.StatefulRefreshTokenRepository {
	return &statefulRefreshTokenRepository{
		db:     db,
		logger: logger,
	}
}

// Save stores a refresh token hash in the repository.
func (r *statefulRefreshTokenRepository) Save(ctx context.Context, tokenHash string, userID string, expiresAt time.Time) error {
	query := `
		INSERT INTO stateful_refresh_tokens (token_hash, user_id, expires_at, revoked)
		VALUES ($1, $2, $3, false)
	`

	_, err := r.db.Exec(ctx, query, tokenHash, userID, expiresAt)
	if err != nil {
		r.logger.Error().Err(err).Str("user_id", userID).Msg("Failed to save stateful refresh token")
		return err
	}

	r.logger.Debug().Str("user_id", userID).Msg("Stateful refresh token saved successfully")
	return nil
}

// IsValid checks if a refresh token is valid and returns the associated user.
func (r *statefulRefreshTokenRepository) IsValid(ctx context.Context, tokenHash string) (*domain.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.email_verified_at, u.password_hash, u.disabled_at, u.created_at, u.updated_at
		FROM stateful_refresh_tokens rt
		JOIN users u ON rt.user_id = u.id
		WHERE rt.token_hash = $1 AND rt.revoked = false AND rt.expires_at > NOW()
	`

	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.EmailVerifiedAt,
		&user.PasswordHash,
		&user.DisabledAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Debug().Msg("Refresh token not found or invalid")
			return nil, ports.ErrTokenNotFound
		}
		r.logger.Error().Err(err).Msg("Error checking refresh token validity")
		return nil, err
	}

	return user, nil
}

// Revoke invalidates a specific refresh token.
func (r *statefulRefreshTokenRepository) Revoke(ctx context.Context, tokenHash string) error {
	query := `UPDATE stateful_refresh_tokens SET revoked = true WHERE token_hash = $1`

	result, err := r.db.Exec(ctx, query, tokenHash)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to revoke refresh token")
		return err
	}

	if result.RowsAffected() == 0 {
		return ports.ErrTokenNotFound
	}

	r.logger.Debug().Msg("Refresh token revoked successfully")
	return nil
}

// RevokeAllForUser invalidates all refresh tokens for a specific user.
func (r *statefulRefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID string) error {
	query := `UPDATE stateful_refresh_tokens SET revoked = true WHERE user_id = $1 AND revoked = false`

	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		r.logger.Error().Err(err).Str("user_id", userID).Msg("Failed to revoke all refresh tokens for user")
		return err
	}

	r.logger.Debug().Str("user_id", userID).Msg("All refresh tokens revoked for user")
	return nil
}
