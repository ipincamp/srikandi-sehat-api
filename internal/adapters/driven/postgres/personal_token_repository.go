package postgres

import (
	"context"

	"github.com/ipincamp/srikandi-sehat/internal/adapters/driven/postgres/models"
	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// Ensure implementation matches the interface
var _ ports.PersonalTokenRepository = (*personalTokenRepository)(nil)

type personalTokenRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

// NewPersonalTokenRepository is the constructor.
func NewPersonalTokenRepository(db *gorm.DB, logger zerolog.Logger) ports.PersonalTokenRepository {
	return &personalTokenRepository{
		db:     db,
		logger: logger,
	}
}

// Save inserts a new token JTI into the database.
func (r *personalTokenRepository) Save(ctx context.Context, token *domain.PersonalToken) error {
	model := models.PersonalTokenFromDomain(token)
	return r.db.WithContext(ctx).Create(model).Error
}

// FindByID finds a token by its JTI (its primary key).
func (r *personalTokenRepository) FindByID(ctx context.Context, jti string) (*domain.PersonalToken, error) {
	var model models.PersonalTokenDBModel
	if err := r.db.WithContext(ctx).Where("id = ?", jti).First(&model).Error; err != nil {
		return nil, err // This will be gorm.ErrRecordNotFound if not found
	}
	return model.ToDomain(), nil
}

// Delete removes a token JTI from the database (logging out).
func (r *personalTokenRepository) Delete(ctx context.Context, jti string) error {
	// We delete using the JTI, which is the 'ID' field
	return r.db.WithContext(ctx).Delete(&models.PersonalTokenDBModel{ID: jti}).Error
}
