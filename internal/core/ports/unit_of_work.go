package ports

import "context"

// TransactionalUnit defines the interface for a single atomic transaction.
// It provides access to repositories that are bound to this specific transaction.
type TransactionalUnit interface {
	// GetUserRepository returns a repository instance bound to this transaction.
	GetUserRepository() UserRepository

	// TODO: Add other transactional repositories here as needed.
	// e.g., GetProductRepository() ProductRepository

	// Commit finalizes all changes made within this transaction.
	Commit(ctx context.Context) error

	// Rollback discards all changes made within this transaction.
	Rollback(ctx context.Context) error
}

// UnitOfWork defines the port for creating and managing transactional units.
// The service layer will depend on this interface to start transactions.
type UnitOfWork interface {
	// Begin starts a new database transaction and returns a TransactionalUnit.
	Begin(ctx context.Context) (TransactionalUnit, error)
}
