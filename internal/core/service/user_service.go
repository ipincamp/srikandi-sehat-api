package service

import (
	"context"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
	"github.com/rs/zerolog"
)

var _ ports.UserService = (*userService)(nil)

// userService adalah implementasi dari ports.UserService
type userService struct {
	repo   ports.UserRepository
	logger zerolog.Logger
	// Nanti Anda akan tambahkan 'password.Hasher' di sini
}

// NewUserService adalah constructor yang benar.
func NewUserService(repo ports.UserRepository, logger zerolog.Logger) ports.UserService {
	return &userService{
		repo:   repo,
		logger: logger,
	}
}

// GetByID implementasi untuk mengambil user
func (s *userService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	s.logger.Info().Str("user_id", id).Msg("Fetching user by ID")

	// Memanggil port repository
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Str("user_id", id).Msg("Failed to find user")
		// TODO: Konversi error (misal: gorm.ErrRecordNotFound ke domain.ErrUserNotFound)
		return nil, err
	}

	return user, nil
}

// UpdateProfile implements the logic for section 1.11.4.
func (s *userService) UpdateProfile(ctx context.Context, userID string, newName string) (*domain.User, error) {
	log := s.logger.With().Str("method", "UpdateProfile").Str("user_id", userID).Logger()

	// 1. Fetch the existing user
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find user for update")
		return nil, err // This will handle "not found" or other DB errors
	}

	// 2. Apply the change
	user.Name = newName

	// 3. Save the updated user object
	// The repository's Update method handles updating the "updated_at" timestamp
	if err := s.repo.Update(ctx, user); err != nil {
		log.Error().Err(err).Msg("Failed to save user updates to repository")
		return nil, err
	}

	log.Info().Msg("User profile updated successfully")

	// 4. Return the updated user, as required
	return user, nil
}
