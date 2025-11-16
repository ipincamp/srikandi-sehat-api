package postgres

import (
	"context"

	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// Ensure implementation matches interfaces
var _ ports.UnitOfWork = (*unitOfWork)(nil)
var _ ports.Transaction = (*transaction)(nil)

// unitOfWork is the GORM implementation of ports.UnitOfWork
type unitOfWork struct {
	db     *gorm.DB
	logger zerolog.Logger
}

// NewUnitOfWork creates a new UoW provider.
func NewUnitOfWork(db *gorm.DB, logger zerolog.Logger) ports.UnitOfWork {
	return &unitOfWork{
		db:     db,
		logger: logger,
	}
}

// Begin starts a new GORM transaction.
func (u *unitOfWork) Begin(ctx context.Context) (ports.Transaction, error) {
	tx := u.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &transaction{
		tx:     tx,
		logger: u.logger,
	}, nil
}

// transaction is the GORM implementation of ports.Transaction
type transaction struct {
	tx     *gorm.DB
	logger zerolog.Logger
}

// GetUserRepository vends a repo bound to the transaction.
func (t *transaction) GetUserRepository() ports.UserRepository {
	// Pass the transaction (t.tx) to the repository constructor
	return NewUserRepository(t.tx, t.logger)
}

// GetPersonalTokenRepository vends a repo bound to the transaction.
func (t *transaction) GetPersonalTokenRepository() ports.PersonalTokenRepository {
	// Pass the transaction (t.tx) to the repository constructor
	return NewPersonalTokenRepository(t.tx, t.logger)
}

// GetUserTokenRepository vends a repo bound to the transaction.
func (t *transaction) GetUserTokenRepository() ports.UserTokenRepository {
	// Pass the transaction (t.tx) to the repository constructor
	return NewUserTokenRepository(t.tx, t.logger)
}

// Commit commits the GORM transaction.
func (t *transaction) Commit() error {
	return t.tx.Commit().Error
}

// Rollback rolls back the GORM transaction.
func (t *transaction) Rollback() error {
	return t.tx.Rollback().Error
}
