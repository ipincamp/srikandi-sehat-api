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

// dbExecutor is an internal interface that abstracts database operations.
// Both *pgxpool.Pool and pgx.Tx satisfy this interface, allowing our
// repositories to be agnostic of whether they are in a transaction or not.
type dbExecutor interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

// Compile-time checks to ensure our types implement the core ports.
var _ ports.UnitOfWork = (*unitOfWork)(nil)
var _ ports.TransactionalUnit = (*transactionalUnit)(nil)

// unitOfWork is the concrete implementation of the ports.UnitOfWork interface.
// It holds the connection pool to begin new transactions.
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

// Begin starts a new pgx transaction and wraps it in our transactionalUnit.
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

// transactionalUnit is the concrete implementation of ports.TransactionalUnit.
// It holds the live database transaction (pgx.Tx).
type transactionalUnit struct {
	tx     pgx.Tx
	logger zerolog.Logger
}

// GetUserRepository creates a new UserRepository instance that is
// bound to this specific transaction (t.tx).
func (t *transactionalUnit) GetUserRepository() ports.UserRepository {
	// Create a new logger context for this repo
	repoLogger := t.logger.With().Str("component", "UserRepository(TX)").Logger()

	// Pass the transaction 't.tx' as the dbExecutor
	// This ensures all operations on the returned repo are part of this transaction.
	return NewUserRepository(t.tx, repoLogger)
}

// Commit commits the underlying pgx transaction.
func (t *transactionalUnit) Commit(ctx context.Context) error {
	if err := t.tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// Rollback rolls back the underlying pgx transaction.
func (t *transactionalUnit) Rollback(ctx context.Context) error {
	if err := t.tx.Rollback(ctx); err != nil {
		// We only log the error here, as a rollback failure
		// is usually secondary to the error that caused the rollback.
		t.logger.Warn().Err(err).Msg("Failed to rollback transaction")
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}
	return nil
}
