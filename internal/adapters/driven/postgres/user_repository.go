package postgres

import (
	"context"
	"time"

	"github.com/ipincamp/srikandi-sehat/internal/adapters/driven/postgres/models"
	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

var _ ports.UserRepository = (*userRepository)(nil)

type userRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

// NewUserRepository adalah constructor untuk DI (Dependency Injection)
func NewUserRepository(db *gorm.DB, logger zerolog.Logger) ports.UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

// FindByID mengambil user dari DB berdasarkan ID
func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var userModel models.UserDBModel

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&userModel).Error; err != nil {
		// Error bisa jadi gorm.ErrRecordNotFound
		return nil, err
	}

	// Konversi model DB ke model Domain murni
	return userModel.ToDomain(), nil
}

// FindByEmail mengambil user dari DB berdasarkan Email
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var userModel models.UserDBModel

	// Gunakan index unik (active email) yang ada di migrasi
	if err := r.db.WithContext(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&userModel).Error; err != nil {
		return nil, err
	}

	return userModel.ToDomain(), nil
}

// Save (untuk membuat user baru)
func (r *userRepository) Save(ctx context.Context, user *domain.User) error {
	userModel := models.FromDomain(user)
	return r.db.WithContext(ctx).Create(userModel).Error
}

// --- Implementasi fungsi lain (Update, Delete, FindAll) ---
func (r *userRepository) FindAll(ctx context.Context) ([]*domain.User, error) {
	panic("not implemented")
}

// Update saves all fields of the domain user struct.
// GORM's .Save() will update all fields for a record if a primary key is present.
func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	// 1. Create a map of fields to update.
	//    This is much safer than .Save() as it's explicit.
	updates := map[string]interface{}{
		"name":              user.Name,
		"email":             user.Email,
		"email_verified_at": user.EmailVerifiedAt,
		"password_hash":     user.PasswordHash,
		"disabled_at":       user.DisabledAt,
		"updated_at":        time.Now(), // Explicitly set UpdatedAt
		// We deliberately omit "created_at" and "deleted_at"
	}

	// 2. Use Model().Where().Updates() for a targeted update.
	//    This targets the user by ID and applies *only* the fields in the map.
	return r.db.WithContext(ctx).
		Model(&models.UserDBModel{}).
		Where("id = ?", user.ID).
		Updates(updates).Error
}

// Delete performs a soft delete on the user.
// GORM will automatically use the "deleted_at" column because
// it is defined as gorm.DeletedAt in the UserDBModel.
func (r *userRepository) Delete(ctx context.Context, id string) error {
	// We execute a GORM Delete operation, targeting the model by its ID.
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.UserDBModel{}).Error
}
