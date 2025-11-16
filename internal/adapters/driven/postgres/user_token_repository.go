package postgres

import (
	"context"

	"github.com/ipincamp/srikandi-sehat/internal/adapters/driven/postgres/models"
	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// Ensure implementation matches the interface.
var _ ports.UserTokenRepository = (*userTokenRepository)(nil)

type userTokenRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

// NewUserTokenRepository is the constructor.
func NewUserTokenRepository(db *gorm.DB, logger zerolog.Logger) ports.UserTokenRepository {
	return &userTokenRepository{
		db:     db,
		logger: logger,
	}
}

// Save uses GORM's Save() for an "Upsert" behavior based on the composite primary key.
// This fulfills the "Update-or-Insert" requirement.
func (r *userTokenRepository) Save(ctx context.Context, token *domain.UserToken) error {
	model := models.UserTokenFromDomain(token)
	return r.db.WithContext(ctx).Save(model).Error
}

// FindByUserIDAndPurpose finds the token in the database.
func (r *userTokenRepository) FindByUserIDAndPurpose(ctx context.Context, userID string, purpose string) (*domain.UserToken, error) {
	var model models.UserTokenDBModel
	if err := r.db.WithContext(ctx).Where("user_id = ? AND purpose = ?", userID, purpose).First(&model).Error; err != nil {
		return nil, err // This will be gorm.ErrRecordNotFound if not found
	}
	return model.ToDomain(), nil
}

// DeleteByUserIDAndPurpose removes one or more tokens.
func (r *userTokenRepository) DeleteByUserIDAndPurpose(ctx context.Context, userID string, purpose string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND purpose = ?", userID, purpose).
		Delete(&models.UserTokenDBModel{}).Error
}
