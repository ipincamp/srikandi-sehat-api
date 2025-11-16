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
