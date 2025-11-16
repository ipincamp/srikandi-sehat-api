package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
	"github.com/ipincamp/srikandi-sehat/internal/core/service"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// setupUserService initializes a UserService with a mock repository.
func setupUserService(t *testing.T) (
	*MockUserRepository, // The mock repo to control
	ports.UserService, // The service to test
) {
	// Create a disabled logger for tests to avoid noisy output.
	logger := zerolog.Nop()

	// Create the mock repository.
	mockRepo := &MockUserRepository{}

	// Initialize the service, injecting the mock.
	userService := service.NewUserService(mockRepo, logger)

	// Ensure the service was created.
	require.NotNil(t, userService)

	return mockRepo, userService
}

func TestGetByID_Success(t *testing.T) {
	// 1. Setup
	mockRepo, userService := setupUserService(t)
	ctx := context.Background()
	testID := "user-uuid-123"

	// Define the expected user to be returned by the mock.
	expectedUser := &domain.User{
		ID:    testID,
		Name:  "Test User",
		Email: "test@example.com",
	}

	// 2. Stub the mock repository method
	// We program our mock to return the expectedUser when FindByID is called.
	mockRepo.FindByIDFunc = func(c context.Context, id string) (*domain.User, error) {
		// We can add assertions here to check if the correct params were passed.
		assert.Equal(t, testID, id, "FindByID called with wrong ID")
		return expectedUser, nil
	}

	// 3. Act
	// Call the service method we are testing.
	user, err := userService.GetByID(ctx, testID)

	// 4. Assert
	// Check that the results are what we expected.
	require.NoError(t, err, "GetByID should not return an error")
	require.NotNil(t, user, "User should not be nil")
	assert.Equal(t, expectedUser.ID, user.ID, "Returned user ID mismatch")
	assert.Equal(t, expectedUser.Name, user.Name, "Returned user Name mismatch")
}

func TestGetByID_NotFound(t *testing.T) {
	// 1. Setup
	mockRepo, userService := setupUserService(t)
	ctx := context.Background()
	testID := "non-existent-id"

	// 2. Stub the mock repository
	// This time, we program the mock to return a "not found" error.
	mockRepo.FindByIDFunc = func(c context.Context, id string) (*domain.User, error) {
		assert.Equal(t, testID, id)
		return nil, gorm.ErrRecordNotFound // Use a real DB error for realism
	}

	// 3. Act
	user, err := userService.GetByID(ctx, testID)

	// 4. Assert
	require.Error(t, err, "GetByID should return an error")
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), "Error should be RecordNotFound")
	assert.Nil(t, user, "User should be nil on error")
}
