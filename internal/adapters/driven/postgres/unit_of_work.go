package postgres

import (
	"context"
	"fmt"

	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

// dbExecutor abstracts database operations allowing repositories to work with both connection pools and transactions.
type dbExecutor interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

// Compile-time interface implementation checks.
var _ ports.UnitOfWork = (*unitOfWork)(nil)
var _ ports.TransactionalUnit = (*transactionalUnit)(nil)

// unitOfWork implements ports.UnitOfWork interface using PostgreSQL connection pool.
type unitOfWork struct {
	pool   *pgxpool.Pool
	logger zerolog.Logger
}

// NewUnitOfWork creates a new UnitOfWork implementation.
func NewUnitOfWork(pool *pgxpool.Pool, logger zerolog.Logger) ports.UnitOfWork {
	return &unitOfWork{
		pool:   pool,
		logger: logger,
	}
}

// Begin starts a new database transaction and returns a transactional unit.
func (u *unitOfWork) Begin(ctx context.Context) (ports.TransactionalUnit, error) {
	// Begin a new transaction from the pool
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Return the transactional unit that holds this transaction
	return &transactionalUnit{
		tx:     tx,
		logger: u.logger,
	}, nil
}

// transactionalUnit implements ports.TransactionalUnit interface with an active database transaction.
type transactionalUnit struct {
	tx     pgx.Tx
	logger zerolog.Logger
}

// GetUserRepository returns a transactional user repository instance.
func (t *transactionalUnit) GetUserRepository() ports.UserRepository {
	// Create a new logger context for this repo
	repoLogger := t.logger.With().Str("component", "UserRepository(TX)").Logger()

	// Pass the transaction 't.tx' as the dbExecutor
	// This ensures all operations on the returned repo are part of this transaction.
	return NewUserRepository(t.tx, repoLogger)
}

// GetUserTokenRepository returns a transactional user token repository instance.
func (t *transactionalUnit) GetUserTokenRepository() ports.UserTokenRepository {
	repoLogger := t.logger.With().Str("component", "UserTokenRepository(TX)").Logger()
	return NewUserTokenRepository(t.tx, repoLogger)
}

// GetRoleRepository returns a transactional role repository instance.
func (t *transactionalUnit) GetRoleRepository() ports.RoleRepository {
	repoLogger := t.logger.With().Str("component", "RoleRepository(TX)").Logger()
	return NewRoleRepository(t.tx, repoLogger)
}

// GetPermissionRepository returns a transactional permission repository instance.
func (t *transactionalUnit) GetPermissionRepository() ports.PermissionRepository {
	repoLogger := t.logger.With().Str("component", "PermissionRepository(TX)").Logger()
	return NewPermissionRepository(t.tx, repoLogger)
}

// GetActivityLogRepository returns a transactional activity log repository instance.
func (t *transactionalUnit) GetActivityLogRepository() ports.ActivityLogRepository {
	repoLogger := t.logger.With().Str("component", "ActivityLogRepository(TX)").Logger()
	return NewActivityLogRepository(t.tx, repoLogger)
}

// GetStatefulRefreshTokenRepository returns a transactional stateful refresh token repository instance.
func (t *transactionalUnit) GetStatefulRefreshTokenRepository() ports.StatefulRefreshTokenRepository {
	repoLogger := t.logger.With().Str("component", "StatefulRefreshTokenRepository(TX)").Logger()
	return NewStatefulRefreshTokenRepository(t.tx, repoLogger)
}

// GetPersonalTokenRepository returns a transactional personal token repository instance.
func (t *transactionalUnit) GetPersonalTokenRepository() ports.PersonalTokenRepository {
	repoLogger := t.logger.With().Str("component", "PersonalTokenRepository(TX)").Logger()
	return NewPersonalTokenRepository(t.tx, repoLogger)
}

// Commit persists all changes made within the transaction to the database.
func (t *transactionalUnit) Commit(ctx context.Context) error {
	if err := t.tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// Rollback discards all changes made within the transaction.
func (t *transactionalUnit) Rollback(ctx context.Context) error {
	if err := t.tx.Rollback(ctx); err != nil {
		t.logger.Warn().Err(err).Msg("Failed to rollback transaction")
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}
	return nil
}
