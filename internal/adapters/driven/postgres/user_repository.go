package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
)

// Compile-time check to ensure userRepository implements ports.UserRepository
var _ ports.UserRepository = (*userRepository)(nil)

// Error specifically for this repository
var ErrUserNotFound = errors.New("user not found")

// userRepository implements the ports.UserRepository interface
// using a pgxpool.Pool for database connections.
type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new repository instance.
func NewUserRepository(db *pgxpool.Pool) ports.UserRepository {
	return &userRepository{db: db}
}

// Save creates a new user in the database.
// It assumes the domain.User has ID=0 for creation.
// Note: This implementation only handles *creation* as requested by AuthService.
// A more robust `Save` would handle updates (if ID > 0).
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
		// TODO: Check for unique constraint violation (pgconn.PgError)
		return err
	}

	// Update the original domain model with new data (ID, timestamps)
	user.ID = dbUser.ID
	user.CreatedAt = dbUser.CreatedAt
	user.UpdatedAt = dbUser.UpdatedAt

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
			return nil, ErrUserNotFound
		}
		return nil, err
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
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return dbUser.toDomain(), nil
}
