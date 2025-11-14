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
var _ ports.UserTokenRepository = (*userTokenRepository)(nil)

// userTokenRepository implements ports.UserTokenRepository interface using PostgreSQL.
type userTokenRepository struct {
	db     dbExecutor
	logger zerolog.Logger
}

// NewUserTokenRepository creates a new user token repository instance.
func NewUserTokenRepository(db dbExecutor, logger zerolog.Logger) ports.UserTokenRepository {
	return &userTokenRepository{
		db:     db,
		logger: logger,
	}
}

// Save persists a new user token to the database.
func (r *userTokenRepository) Save(ctx context.Context, token *domain.UserToken) error {
	now := time.Now().UTC()

	query := `
		INSERT INTO user_tokens (user_id, purpose, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, query,
		token.UserID,
		token.Purpose,
		token.TokenHash,
		token.ExpiresAt,
		now,
	)

	if err != nil {
		r.logger.Error().Err(err).Str("user_id", token.UserID).Str("purpose", token.Purpose).Msg("Failed to save user token")
		return err
	}

	token.CreatedAt = now
	r.logger.Debug().Str("user_id", token.UserID).Str("purpose", token.Purpose).Msg("User token saved successfully")
	return nil
}

// FindAndConsume retrieves and marks a token as consumed in an atomic operation.
func (r *userTokenRepository) FindAndConsume(ctx context.Context, userID, purpose, tokenHash string) (*domain.UserToken, error) {
	query := `
		DELETE FROM user_tokens
		WHERE user_id = $1 AND purpose = $2 AND token_hash = $3 AND expires_at > NOW()
		RETURNING user_id, purpose, token_hash, expires_at, created_at
	`

	token := &domain.UserToken{}
	err := r.db.QueryRow(ctx, query, userID, purpose, tokenHash).Scan(
		&token.UserID,
		&token.Purpose,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Debug().Str("user_id", userID).Str("purpose", purpose).Msg("User token not found or expired")
			return nil, ports.ErrTokenNotFound
		}
		r.logger.Error().Err(err).Str("user_id", userID).Str("purpose", purpose).Msg("Error finding and consuming user token")
		return nil, err
	}

	r.logger.Debug().Str("user_id", userID).Str("purpose", purpose).Msg("User token found and consumed")
	return token, nil
}
