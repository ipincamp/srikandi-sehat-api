package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
)

// Compile-time interface implementation check.
var _ ports.PersonalTokenRepository = (*personalTokenRepository)(nil)

// personalTokenRepository implements ports.PersonalTokenRepository interface using PostgreSQL.
type personalTokenRepository struct {
	db     dbExecutor
	logger zerolog.Logger
}

// NewPersonalTokenRepository creates a new personal token repository instance.
func NewPersonalTokenRepository(db dbExecutor, logger zerolog.Logger) ports.PersonalTokenRepository {
	return &personalTokenRepository{
		db:     db,
		logger: logger,
	}
}

// Save persists a new personal access token to the repository.
func (r *personalTokenRepository) Save(ctx context.Context, token *domain.PersonalToken) error {
	query := `
		INSERT INTO personal_tokens (id, user_id, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET user_id = $2, expires_at = $3
	`

	_, err := r.db.Exec(ctx, query, token.ID, token.UserID, token.ExpiresAt)
	if err != nil {
		r.logger.Error().Err(err).Str("token_id", token.ID).Msg("Failed to save personal token")
		return err
	}

	r.logger.Debug().Str("token_id", token.ID).Msg("Personal token saved successfully")
	return nil
}

// FindByID retrieves a personal access token by its unique identifier.
func (r *personalTokenRepository) FindByID(ctx context.Context, tokenID string) (*domain.PersonalToken, error) {
	query := `SELECT id, user_id, expires_at FROM personal_tokens WHERE id = $1`

	token := &domain.PersonalToken{}
	err := r.db.QueryRow(ctx, query, tokenID).Scan(&token.ID, &token.UserID, &token.ExpiresAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Debug().Str("token_id", tokenID).Msg("Personal token not found")
			return nil, ports.ErrTokenNotFound
		}
		r.logger.Error().Err(err).Str("token_id", tokenID).Msg("Error finding personal token by ID")
		return nil, err
	}

	return token, nil
}

// ListByUserID retrieves all personal access tokens belonging to a user.
func (r *personalTokenRepository) ListByUserID(ctx context.Context, userID string) ([]*domain.PersonalToken, error) {
	query := `SELECT id, user_id, expires_at FROM personal_tokens WHERE user_id = $1 ORDER BY expires_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		r.logger.Error().Err(err).Str("user_id", userID).Msg("Failed to list personal tokens by user ID")
		return nil, err
	}
	defer rows.Close()

	var tokens []*domain.PersonalToken
	for rows.Next() {
		token := &domain.PersonalToken{}
		if err := rows.Scan(&token.ID, &token.UserID, &token.ExpiresAt); err != nil {
			r.logger.Error().Err(err).Msg("Error scanning personal token")
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, rows.Err()
}

// Delete removes a personal access token from the repository.
func (r *personalTokenRepository) Delete(ctx context.Context, tokenID string, userID string) error {
	query := `DELETE FROM personal_tokens WHERE id = $1 AND user_id = $2`

	result, err := r.db.Exec(ctx, query, tokenID, userID)
	if err != nil {
		r.logger.Error().Err(err).Str("token_id", tokenID).Str("user_id", userID).Msg("Failed to delete personal token")
		return err
	}

	if result.RowsAffected() == 0 {
		return ports.ErrTokenNotFound
	}

	r.logger.Debug().Str("token_id", tokenID).Str("user_id", userID).Msg("Personal token deleted successfully")
	return nil
}
