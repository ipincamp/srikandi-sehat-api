package ports

import "context"

// TransactionalUnit represents a database transaction that provides access to transactional repositories.
type TransactionalUnit interface {
	// Commit persists all changes made within the transaction to the database.
	Commit(ctx context.Context) error

	// Rollback discards all changes made within the transaction.
	Rollback(ctx context.Context) error

	// GetUserRepository returns a transactional instance of UserRepository.
	GetUserRepository() UserRepository

	// GetUserTokenRepository returns a transactional instance of UserTokenRepository.
	GetUserTokenRepository() UserTokenRepository

	// GetRoleRepository returns a transactional instance of RoleRepository.
	GetRoleRepository() RoleRepository

	// GetPermissionRepository returns a transactional instance of PermissionRepository.
	GetPermissionRepository() PermissionRepository

	// GetActivityLogRepository returns a transactional instance of ActivityLogRepository.
	GetActivityLogRepository() ActivityLogRepository

	// GetStatefulRefreshTokenRepository returns a transactional instance of StatefulRefreshTokenRepository.
	GetStatefulRefreshTokenRepository() StatefulRefreshTokenRepository

	// GetPersonalTokenRepository returns a transactional instance of PersonalTokenRepository.
	GetPersonalTokenRepository() PersonalTokenRepository
}

// UnitOfWork manages database transactions ensuring atomicity across multiple repository operations.
type UnitOfWork interface {
	// Begin starts a new database transaction and returns a TransactionalUnit.
	Begin(ctx context.Context) (TransactionalUnit, error)
}
