package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/rs/zerolog"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
)

// Compile-time check to ensure userRepository implements ports.UserRepository
var _ ports.UserRepository = (*userRepository)(nil)

// userRepository implements the ports.UserRepository interface
// using a pgxpool.Pool for database connections.
type userRepository struct {
	db     dbExecutor
	logger zerolog.Logger
}

// NewUserRepository creates a new repository instance.
func NewUserRepository(db dbExecutor, logger zerolog.Logger) ports.UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

// Save creates a new user in the database.
func (r *userRepository) Save(ctx context.Context, user *domain.User) error {
	// Map domain model to database model
	dbUser := fromDomain(user)

	// Set timestamps for creation
	now := time.Now().UTC()
	dbUser.CreatedAt = now
	dbUser.UpdatedAt = now

	query := `
		INSERT INTO users (uuid, name, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		dbUser.UUID,
		dbUser.Name,
		dbUser.Email,
		dbUser.Password,
		dbUser.CreatedAt,
		dbUser.UpdatedAt,
	).Scan(&dbUser.ID, &dbUser.CreatedAt, &dbUser.UpdatedAt)

	if err != nil {
		// --- Enhanced Error Handling ---
		var pgErr *pgconn.PgError
		// Check if the error is a PostgreSQL error
		if errors.As(err, &pgErr) {
			// Check for unique_violation (e.g., duplicate email)
			if pgErr.Code == "23505" {
				r.logger.Warn().
					Str("email", dbUser.Email).
					Str("constraint", pgErr.ConstraintName).
					Msg("User save failed: unique constraint violation")
				return ports.ErrDuplicateEmail
			}
		}

		// Log any other database error
		r.logger.Error().Err(err).Str("email", dbUser.Email).Msg("Failed to save user")
		return ports.ErrUnexpectedSave
	}

	// Update the original domain model with new data (ID, timestamps)
	user.ID = dbUser.ID
	user.CreatedAt = dbUser.CreatedAt
	user.UpdatedAt = dbUser.UpdatedAt

	r.logger.Debug().Str("uuid", user.UUID).Msg("User saved successfully")
	return nil
}

// FindByID retrieves a user by their public UUID.
func (r *userRepository) FindByID(ctx context.Context, uuid string) (*domain.User, error) {
	query := `
		SELECT id, uuid, name, email, password, created_at, updated_at
		FROM users
		WHERE uuid = $1
	`

	dbUser := &User{}
	err := r.db.QueryRow(ctx, query, uuid).Scan(
		&dbUser.ID,
		&dbUser.UUID,
		&dbUser.Name,
		&dbUser.Email,
		&dbUser.Password,
		&dbUser.CreatedAt,
		&dbUser.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Debug().Str("uuid", uuid).Msg("User not found by UUID")
			return nil, ports.ErrUserNotFound
		}
		// Log other errors
		r.logger.Error().Err(err).Str("uuid", uuid).Msg("Error finding user by UUID")
		return nil, ports.ErrUnexpectedFind
	}

	return dbUser.toDomain(), nil
}

// FindByEmail retrieves a user by their email address.
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, uuid, name, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	dbUser := &User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&dbUser.ID,
		&dbUser.UUID,
		&dbUser.Name,
		&dbUser.Email,
		&dbUser.Password,
		&dbUser.CreatedAt,
		&dbUser.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Debug().Str("email", email).Msg("User not found by email")
			return nil, ports.ErrUserNotFound
		}
		// Log other errors
		r.logger.Error().Err(err).Str("email", email).Msg("Error finding user by email")
		return nil, ports.ErrUnexpectedFind
	}

	return dbUser.toDomain(), nil
}

// FindMapByUUIDs implements the batch fetch method for the UserRepository port.
func (r *userRepository) FindMapByUUIDs(ctx context.Context, uuids []string) (map[string]*domain.User, error) {
	query := `
		SELECT id, uuid, name, email, password, created_at, updated_at
		FROM users
		WHERE uuid = ANY($1)
	`
	// Handle empty input UUIDs to avoid a query error
	if len(uuids) == 0 {
		return make(map[string]*domain.User), nil
	}

	rows, err := r.db.Query(ctx, query, uuids)
	if err != nil {
		r.logger.Error().Err(err).Strs("uuids", uuids).Msg("Failed to query users by UUIDs")
		return nil, ports.ErrUnexpectedFind
	}
	defer rows.Close()

	// Initialize a map to store results, keyed by UUID
	userMap := make(map[string]*domain.User, len(uuids))

	for rows.Next() {
		dbUser := &User{}
		err := rows.Scan(
			&dbUser.ID,
			&dbUser.UUID,
			&dbUser.Name,
			&dbUser.Email,
			&dbUser.Password,
			&dbUser.CreatedAt,
			&dbUser.UpdatedAt,
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("Error scanning user row in batch find")
			return nil, ports.ErrUnexpectedFind
		}
		// Map the database model to domain model and add to map
		userMap[dbUser.UUID] = dbUser.toDomain()
	}

	if err := rows.Err(); err != nil {
		r.logger.Error().Err(err).Msg("Error after iterating user rows in batch find")
		return nil, ports.ErrUnexpectedFind
	}

	r.logger.Debug().Int("found", len(userMap)).Int("requested", len(uuids)).Msg("Batch user find complete")
	return userMap, nil
}
