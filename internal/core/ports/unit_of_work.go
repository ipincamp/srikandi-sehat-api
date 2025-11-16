package ports

import "context"

// Transaction defines the interface for a single atomic transaction.
// It provides access to transactional repositories.
type Transaction interface {
	// GetUserRepository returns a UserRepository bound to this transaction.
	GetUserRepository() UserRepository
	// GetPersonalTokenRepository returns a PersonalTokenRepository bound to this transaction.
	GetPersonalTokenRepository() PersonalTokenRepository

	// Commit finalizes the transaction.
	Commit() error
	// Rollback discards the transaction.
	Rollback() error
}

// UnitOfWork defines the interface for starting and managing transactions.
type UnitOfWork interface {
	// Begin starts a new transaction.
	Begin(ctx context.Context) (Transaction, error)
}
