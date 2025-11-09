package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
	"github.com/ipincamp/srikandi-sehat/pkg/password"

	db "github.com/ipincamp/srikandi-sehat/internal/adapters/driven/postgres"
)

// Compile-time check
var _ ports.UserService = (*userService)(nil)

// userService implements the ports.UserService interface.
type userService struct {
	userRepo ports.UserRepository // For non-transactional reads
	hasher   password.Hasher
	logger   zerolog.Logger
	uow      ports.UnitOfWork
}

// NewUserService is the constructor for userService.
func NewUserService(
	userRepo ports.UserRepository,
	hasher password.Hasher,
	logger zerolog.Logger,
	uow ports.UnitOfWork,
) ports.UserService {
	return &userService{
		userRepo: userRepo,
		hasher:   hasher,
		logger:   logger,
		uow:      uow,
	}
}

// GetUserByID retrieves a user's public profile information.
func (s *userService) GetUserByID(ctx context.Context, uuid string) (*domain.User, error) {
	// This read operation is non-transactional and can use the base repo.
	user, err := s.userRepo.FindByID(ctx, uuid)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			s.logger.Warn().Str("uuid", uuid).Msg("User not found")
			return nil, ErrUserNotFound // Use service-level error
		}
		s.logger.Error().Err(err).Str("uuid", uuid).Msg("Failed to get user by ID")
		return nil, err
	}
	return user, nil
}

// CreateUser creates a new user (e.g., for an admin panel).
// This is distinct from 'Register' as it does not return tokens.
func (s *userService) CreateUser(ctx context.Context, name, email, passwordStr string) (*domain.User, error) {
	// 1. Pre-check existence using the non-transactional repo
	_, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		// User found, email is taken
		s.logger.Warn().Str("email", email).Msg("CreateUser failed: email already exists (pre-check)")
		return nil, ErrEmailExists
	}
	if !errors.Is(err, db.ErrUserNotFound) {
		// A different, unexpected database error occurred
		s.logger.Error().Err(err).Str("email", email).Msg("Failed to check user existence")
		return nil, err
	}

	// 2. Hash the password
	hashedPassword, err := s.hasher.Hash(passwordStr)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to hash password during user creation")
		return nil, err
	}

	// 3. Create the domain user
	user := &domain.User{
		UUID:     uuid.NewString(),
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	// 4. === Begin Transactional Unit of Work ===
	tx, err := s.uow.Begin(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to begin CreateUser transaction")
		return nil, err
	}

	// Defer rollback in case of error
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				s.logger.Error().Err(rbErr).Msg("Failed to rollback transaction after error")
			}
		}
	}()

	// 5. Get the transactional repository
	txUserRepo := tx.GetUserRepository()

	// 6. Save the user *using the transactional repo*
	if err = txUserRepo.Save(ctx, user); err != nil {
		// Check for duplicate email (race condition)
		if errors.Is(err, db.ErrDuplicateEmail) {
			s.logger.Warn().Str("email", email).Msg("CreateUser failed: email already exists (race condition on save)")
			return nil, ErrEmailExists
		}
		s.logger.Error().Err(err).Str("email", email).Msg("Failed to save user during creation")
		return nil, err // Defer will catch this and rollback
	}

	// 7. Commit the transaction
	if err = tx.Commit(ctx); err != nil {
		s.logger.Error().Err(err).Msg("Failed to commit CreateUser transaction")
		return nil, err
	}
	// === End Transactional Unit of Work ===

	s.logger.Info().Str("email", email).Str("uuid", user.UUID).Msg("User created successfully via CreateUser")

	// Note: We don't add to the bloom filter here, as 'Register' is the primary path.
	// Or, if this is a valid path, we should also inject 'userCache' and call 'userCache.Add(user.Email)'.
	// For now, we follow the 'Register' service's pattern.

	return user, nil
}
