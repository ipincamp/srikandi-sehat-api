package service_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
	"github.com/ipincamp/srikandi-sehat/internal/core/service"
	"github.com/ipincamp/srikandi-sehat/pkg/config"
	"github.com/ipincamp/srikandi-sehat/pkg/token"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// Test complex setup for AuthService
type authServiceMocks struct {
	UOW           *MockUnitOfWork
	Tx            *MockTransaction
	UserRepo      *MockUserRepository
	TokenRepo     *MockPersonalTokenRepository
	UserTokenRepo *MockUserTokenRepository
	Hasher        *MockHasher
	Maker         *MockTokenMaker
	Mail          *MockMailService
	TokenCfg      config.Token
}

// setupAuthService initializes an AuthService with all its dependencies mocked.
func setupAuthService(t *testing.T) (authServiceMocks, ports.AuthService) {
	logger := zerolog.Nop()

	// 1. Create all mock dependencies
	mocks := authServiceMocks{
		UOW:           &MockUnitOfWork{},
		UserRepo:      &MockUserRepository{},
		TokenRepo:     &MockPersonalTokenRepository{},
		UserTokenRepo: &MockUserTokenRepository{},
		Hasher:        &MockHasher{},
		Maker:         &MockTokenMaker{},
		Mail:          &MockMailService{},
		TokenCfg: config.Token{
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 7 * 24 * time.Hour,
		},
	}

	// 2. Create the mock transaction that holds the mock repos
	mocks.Tx = &MockTransaction{
		MockUserRepo:      mocks.UserRepo,
		MockTokenRepo:     mocks.TokenRepo,
		MockUserTokenRepo: mocks.UserTokenRepo,
	}

	// 3. Initialize the service, injecting all mocks
	authService := service.NewAuthService(
		mocks.UOW,
		mocks.Hasher,
		mocks.Maker,
		mocks.TokenCfg,
		logger,
		mocks.Mail,
	)
	require.NotNil(t, authService)

	return mocks, authService
}

func TestRegister_Success(t *testing.T) {
	// 1. Setup
	mocks, authService := setupAuthService(t)
	ctx := context.Background()

	// Test inputs
	name := "New User"
	email := "new@example.com"
	password := "Password123!"

	// Mocked data
	hashedPassword := "hashed_password"
	newUserID := uuid.NewString()
	accessToken := "new_access_token"
	refreshToken := "new_refresh_token"
	jti := uuid.NewString()

	// 2. Stub all mock calls in order of execution

	// a. UOW begins transaction
	mocks.UOW.BeginFunc = func(c context.Context) (ports.Transaction, error) {
		return mocks.Tx, nil
	}

	// b. UserRepo checks if email exists (it doesn't)
	mocks.UserRepo.FindByEmailFunc = func(c context.Context, e string) (*domain.User, error) {
		assert.Equal(t, email, e)

		// Use the struct field for state, not context
		if !mocks.UserRepo.FindByEmailCalled {
			// First call: checking for existence
			mocks.UserRepo.FindByEmailCalled = true // Mutate the mock's state
			return nil, gorm.ErrRecordNotFound
		}
		// Second call: fetching the *newly created* user
		return &domain.User{ID: newUserID, Name: name, Email: email}, nil
	}

	// c. Hasher hashes the password
	mocks.Hasher.HashFunc = func(p string) (string, error) {
		assert.Equal(t, password, p)
		return hashedPassword, nil
	}

	// d. UserRepo saves the new user
	mocks.UserRepo.SaveFunc = func(c context.Context, u *domain.User) error {
		assert.Equal(t, name, u.Name)
		assert.Equal(t, email, u.Email)
		assert.Equal(t, hashedPassword, u.PasswordHash)
		return nil
	}

	// e. TokenMaker creates access token
	mocks.Maker.CreateTokenFunc = func(userID, roleID, useFor string, duration time.Duration) (string, *token.Payload, error) {
		if useFor == token.UseForAccessToken {
			assert.Equal(t, newUserID, userID)
			return accessToken, &token.Payload{UserID: newUserID}, nil
		}
		// f. TokenMaker creates refresh token
		if useFor == token.UseForRefreshToken {
			assert.Equal(t, newUserID, userID)
			return refreshToken, &token.Payload{UserID: newUserID, JTI: jti}, nil
		}
		return "", nil, errors.New("unexpected token type")
	}

	// g. TokenRepo saves the new JTI
	mocks.TokenRepo.SaveFunc = func(c context.Context, pt *domain.PersonalToken) error {
		assert.Equal(t, jti, pt.ID)
		assert.Equal(t, newUserID, pt.UserID)
		return nil
	}

	// h. Transaction commits
	mocks.Tx.CommitFunc = func() error {
		return nil
	}

	// i. Transaction rolls back
	//    This stub will only be called if the test fails unexpectedly.
	mocks.Tx.RollbackFunc = func() error {
		t.Log("Unexpected rollback called") // Log this so we know something is wrong
		return nil
	}

	// j. MailService sends email (asynchronously, stub so it doesn't error)
	mocks.Mail.SendFunc = func(c context.Context, to, subject, plainBody, htmlBody string) error {
		assert.Equal(t, email, to)
		return nil
	}

	// 3. Act
	resp, err := authService.Register(ctx, name, email, password)

	// 4. Assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, accessToken, resp.AccessToken)
	assert.Equal(t, refreshToken, resp.RefreshToken)
}

func TestRegister_EmailAlreadyInUse(t *testing.T) {
	// 1. Setup
	mocks, authService := setupAuthService(t)
	ctx := context.Background()

	// 2. Stub Mocks
	// a. UOW begins transaction
	mocks.UOW.BeginFunc = func(c context.Context) (ports.Transaction, error) {
		return mocks.Tx, nil
	}

	// b. UserRepo finds an *existing* user
	mocks.UserRepo.FindByEmailFunc = func(c context.Context, e string) (*domain.User, error) {
		return &domain.User{ID: "existing-id", Email: e}, nil // Success, user found
	}

	// c. Transaction rolls back
	mocks.Tx.RollbackFunc = func() error {
		return nil
	}

	// 3. Act
	resp, err := authService.Register(ctx, "Test", "existing@example.com", "password")

	// 4. Assert
	require.Error(t, err)
	assert.Equal(t, "email already in use", err.Error())
	assert.Nil(t, resp)
}

func TestLogin_Success(t *testing.T) {
	// 1. Setup
	mocks, authService := setupAuthService(t)
	ctx := context.Background()

	email := "user@example.com"
	password := "Password123!"
	hashedPassword := "hashed_password"
	userID := uuid.NewString()

	// 2. Stub Mocks
	// a. UOW begins
	mocks.UOW.BeginFunc = func(c context.Context) (ports.Transaction, error) {
		return mocks.Tx, nil
	}

	// b. UserRepo finds the user
	mocks.UserRepo.FindByEmailFunc = func(c context.Context, e string) (*domain.User, error) {
		assert.Equal(t, email, e)
		return &domain.User{ID: userID, Email: email, PasswordHash: hashedPassword}, nil
	}

	// c. Hasher compares password (success)
	mocks.Hasher.CompareFunc = func(h string, p string) bool {
		assert.Equal(t, hashedPassword, h)
		assert.Equal(t, password, p)
		return true // Passwords match
	}

	// d. TokenMaker creates tokens (simplified)
	mocks.Maker.CreateTokenFunc = func(uid, rid, useFor string, dur time.Duration) (string, *token.Payload, error) {
		if useFor == token.UseForAccessToken {
			return "access_token", &token.Payload{UserID: uid}, nil
		}
		return "refresh_token", &token.Payload{UserID: uid, JTI: "jti-123"}, nil
	}

	// e. TokenRepo saves JTI
	mocks.TokenRepo.SaveFunc = func(c context.Context, pt *domain.PersonalToken) error {
		return nil
	}

	// f. Transaction commits
	mocks.Tx.CommitFunc = func() error {
		return nil
	}

	// 3. Act
	resp, err := authService.Login(ctx, email, password)

	// 4. Assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "access_token", resp.AccessToken)
}

func TestLogin_WrongPassword(t *testing.T) {
	// 1. Setup
	mocks, authService := setupAuthService(t)
	ctx := context.Background()

	// 2. Stub Mocks
	// a. UOW begins
	mocks.UOW.BeginFunc = func(c context.Context) (ports.Transaction, error) {
		return mocks.Tx, nil
	}

	// b. UserRepo finds user
	mocks.UserRepo.FindByEmailFunc = func(c context.Context, e string) (*domain.User, error) {
		return &domain.User{ID: "id-123", Email: e, PasswordHash: "real_hash"}, nil
	}

	// c. Hasher compares password (fail)
	mocks.Hasher.CompareFunc = func(h string, p string) bool {
		return false // Passwords do NOT match
	}

	// d. Transaction rolls back
	mocks.Tx.RollbackFunc = func() error {
		return nil
	}

	// 3. Act
	resp, err := authService.Login(ctx, "user@example.com", "wrong_password")

	// 4. Assert
	require.Error(t, err)
	assert.Equal(t, "invalid email or password", err.Error())
	assert.Nil(t, resp)
}

func TestRefreshToken_Success(t *testing.T) {
	// 1. Setup
	mocks, authService := setupAuthService(t)
	ctx := context.Background()

	// Mocked data
	oldRefreshToken := "old_refresh_token"
	oldJTI := "old-jti-123"
	userID := "user-uuid-abc"

	newAccessToken := "new_access_token"
	newRefreshToken := "new_refresh_token"
	newJTI := "new-jti-456"

	// 2. Stub all mock calls
	// a. Token validation succeeds (Requirement 1.4.3.1)
	mocks.Maker.ValidateTokenFunc = func(tokenStr string) (*token.Payload, error) {
		assert.Equal(t, oldRefreshToken, tokenStr)
		return &token.Payload{
			JTI:    oldJTI,
			UserID: userID,
			UseFor: token.UseForRefreshToken,
		}, nil
	}

	// b. UOW begins transaction
	mocks.UOW.BeginFunc = func(c context.Context) (ports.Transaction, error) {
		return mocks.Tx, nil
	}

	// c. TokenRepo finds the old JTI in the DB (Requirement 1.4.3.3)
	mocks.TokenRepo.FindByIDFunc = func(c context.Context, jti string) (*domain.PersonalToken, error) {
		assert.Equal(t, oldJTI, jti)
		return &domain.PersonalToken{ID: oldJTI, UserID: userID}, nil
	}

	// d. TokenRepo deletes the old JTI (Requirement 1.4.3.8.1)
	mocks.TokenRepo.DeleteFunc = func(c context.Context, jti string) error {
		assert.Equal(t, oldJTI, jti)
		return nil
	}

	// e. UserRepo finds the user for the new token (Requirement 1.4.3.6)
	mocks.UserRepo.FindByIDFunc = func(c context.Context, id string) (*domain.User, error) {
		assert.Equal(t, userID, id)
		return &domain.User{ID: userID, Name: "Test User"}, nil
	}

	// f. TokenMaker creates the new tokens (Requirement 1.4.3.8.2 & .3)
	mocks.Maker.CreateTokenFunc = func(uid, rid, useFor string, dur time.Duration) (string, *token.Payload, error) {
		if useFor == token.UseForAccessToken {
			return newAccessToken, &token.Payload{UserID: uid}, nil
		}
		// This is the new refresh token
		return newRefreshToken, &token.Payload{UserID: uid, JTI: newJTI}, nil
	}

	// g. TokenRepo saves the new JTI (Requirement 1.4.3.8.4)
	mocks.TokenRepo.SaveFunc = func(c context.Context, pt *domain.PersonalToken) error {
		assert.Equal(t, newJTI, pt.ID)
		assert.Equal(t, userID, pt.UserID)
		return nil
	}

	// h. Transaction commits
	mocks.Tx.CommitFunc = func() error {
		return nil
	}

	// 3. Act
	resp, err := authService.RefreshToken(ctx, oldRefreshToken)

	// 4. Assert (Requirement 1.4.3.10)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, newAccessToken, resp.AccessToken)
	assert.Equal(t, newRefreshToken, resp.RefreshToken)
}

func TestRefreshToken_JTI_NotFound(t *testing.T) {
	// 1. Setup
	mocks, authService := setupAuthService(t)
	ctx := context.Background()
	oldRefreshToken := "already_used_token"
	oldJTI := "jti-that-was-revoked"

	// 2. Stub Mocks
	// a. Token validation succeeds (the token itself is valid)
	mocks.Maker.ValidateTokenFunc = func(tokenStr string) (*token.Payload, error) {
		return &token.Payload{
			JTI:    oldJTI,
			UserID: "user-uuid-abc",
			UseFor: token.UseForRefreshToken,
		}, nil
	}

	// b. UOW begins transaction
	mocks.UOW.BeginFunc = func(c context.Context) (ports.Transaction, error) {
		return mocks.Tx, nil
	}

	// c. TokenRepo *fails* to find the JTI (Requirement 1.4.3.4)
	mocks.TokenRepo.FindByIDFunc = func(c context.Context, jti string) (*domain.PersonalToken, error) {
		// This simulates a token reuse attack or a revoked token.
		return nil, gorm.ErrRecordNotFound
	}

	// d. Transaction rolls back
	mocks.Tx.RollbackFunc = func() error {
		return nil
	}

	// 3. Act
	resp, err := authService.RefreshToken(ctx, oldRefreshToken)

	// 4. Assert
	require.Error(t, err)
	// This error comes from our service logic
	assert.Equal(t, "invalid refresh token", err.Error())
	assert.Nil(t, resp)
}

func TestLogout_Success(t *testing.T) {
	// 1. Setup
	mocks, authService := setupAuthService(t)
	ctx := context.Background()

	refreshToken := "valid_refresh_token"
	jti := "jti-to-revoke"
	userID := "user-uuid-abc"

	// 2. Stub Mocks
	// a. Token validation succeeds (Requirement 1.3.3.1)
	mocks.Maker.ValidateTokenFunc = func(tokenStr string) (*token.Payload, error) {
		assert.Equal(t, refreshToken, tokenStr)
		return &token.Payload{
			JTI:    jti,
			UserID: userID,
			UseFor: token.UseForRefreshToken,
		}, nil
	}

	// b. UOW begins transaction
	mocks.UOW.BeginFunc = func(c context.Context) (ports.Transaction, error) {
		return mocks.Tx, nil
	}

	// c. TokenRepo finds the JTI in the DB (Requirement 1.3.3.3)
	mocks.TokenRepo.FindByIDFunc = func(c context.Context, j string) (*domain.PersonalToken, error) {
		assert.Equal(t, jti, j)
		return &domain.PersonalToken{ID: jti, UserID: userID}, nil
	}

	// d. TokenRepo deletes the JTI (Requirement 1.3.3.7)
	mocks.TokenRepo.DeleteFunc = func(c context.Context, j string) error {
		assert.Equal(t, jti, j)
		return nil
	}

	// e. Transaction commits
	mocks.Tx.CommitFunc = func() error {
		return nil
	}

	// 3. Act
	err := authService.Logout(ctx, refreshToken)

	// 4. Assert (Requirement 1.3.3.8)
	require.NoError(t, err)
}

func TestLogout_AlreadyLoggedOut(t *testing.T) {
	// 1. Setup
	mocks, authService := setupAuthService(t)
	ctx := context.Background()
	refreshToken := "already_logged_out_token"
	jti := "jti-revoked-earlier"

	// 2. Stub Mocks
	// a. Token validation succeeds
	mocks.Maker.ValidateTokenFunc = func(tokenStr string) (*token.Payload, error) {
		return &token.Payload{
			JTI:    jti,
			UserID: "user-uuid-abc",
			UseFor: token.UseForRefreshToken,
		}, nil
	}

	// b. UOW begins transaction
	mocks.UOW.BeginFunc = func(c context.Context) (ports.Transaction, error) {
		return mocks.Tx, nil
	}

	// c. TokenRepo *fails* to find the JTI (Requirement 1.3.3.4)
	mocks.TokenRepo.FindByIDFunc = func(c context.Context, j string) (*domain.PersonalToken, error) {
		return nil, gorm.ErrRecordNotFound
	}

	// d. Transaction still commits (service logic treats this as success)
	mocks.Tx.CommitFunc = func() error {
		return nil
	}

	// 3. Act
	err := authService.Logout(ctx, refreshToken)

	// 4. Assert
	// The service correctly returns no error, as the token is effectively logged out.
	require.NoError(t, err)
}

// This new test verifies the fix for your scenario
func TestResetPassword_Success(t *testing.T) {
	// 1. Setup
	mocks, authService := setupAuthService(t)
	ctx := context.Background()

	// Test inputs
	userID := "user-uuid-123"
	plainToken := "my-secret-reset-token-string"
	// The token format required by the service is "uid.tokenstring"
	fullToken := fmt.Sprintf("%s.%s", userID, plainToken)
	newPassword := "NewS3cureP@ssword!"

	// Mocked data
	hashedPlainToken := "fe1fba9a0f2a5f75deaa2b136fc87079c2a6475bc85ec4469d1ab903b6d09cfe"
	newHashedPassword := "new-argon-hashed-password"

	// This is the user object as it exists in the DB *before* the update
	userToUpdate := &domain.User{
		ID:           userID,
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "old-argon-hash",
	}

	// This is the token object as it exists in the DB
	resetToken := &domain.UserToken{
		UserID:    userID,
		Purpose:   domain.TokenPurposePasswordReset,
		TokenHash: hashedPlainToken, // This must match the hash of plainToken
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	// 2. Stub Mocks (in order of execution)
	// a. UOW begins transaction
	mocks.UOW.BeginFunc = func(c context.Context) (ports.Transaction, error) {
		return mocks.Tx, nil
	}

	// b. UserTokenRepo finds the reset token in the DB
	mocks.UserTokenRepo.FindByUserIDAndPurposeFunc = func(c context.Context, uid string, purpose string) (*domain.UserToken, error) {
		assert.Equal(t, userID, uid)
		assert.Equal(t, domain.TokenPurposePasswordReset, purpose)
		return resetToken, nil
	}

	// c. Hasher hashes the *new* password
	mocks.Hasher.HashFunc = func(p string) (string, error) {
		assert.Equal(t, newPassword, p)
		return newHashedPassword, nil
	}

	// d. UserRepo finds the user to update
	mocks.UserRepo.FindByIDFunc = func(c context.Context, id string) (*domain.User, error) {
		assert.Equal(t, userID, id)
		return userToUpdate, nil
	}

	// e. UserRepo *updates* the user (This is what we fixed)
	mocks.UserRepo.UpdateFunc = func(c context.Context, u *domain.User) error {
		// Assert that the user object being passed for update
		// has the *new* password hash.
		assert.Equal(t, userID, u.ID)
		assert.Equal(t, newHashedPassword, u.PasswordHash)
		return nil
	}

	// f. PersonalTokenRepo deletes all refresh tokens (log out everywhere)
	mocks.TokenRepo.DeleteByUserIDFunc = func(c context.Context, uid string) error {
		assert.Equal(t, userID, uid)
		return nil
	}

	// g. UserTokenRepo deletes the used reset token
	mocks.UserTokenRepo.DeleteByUserIDAndPurposeFunc = func(c context.Context, uid string, purpose string) error {
		assert.Equal(t, userID, uid)
		assert.Equal(t, domain.TokenPurposePasswordReset, purpose)
		return nil
	}

	// h. Transaction commits
	mocks.Tx.CommitFunc = func() error {
		return nil
	}

	// i. Transaction rolls back (should not be called)
	mocks.Tx.RollbackFunc = func() error {
		t.Log("Unexpected rollback called")
		return nil
	}

	// 3. Act
	err := authService.ResetPassword(ctx, fullToken, newPassword)

	// 4. Assert
	require.NoError(t, err)
}
