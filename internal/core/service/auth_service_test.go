package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
	"github.com/ipincamp/srikandi-sehat/internal/core/service"
	"github.com/ipincamp/srikandi-sehat/pkg/config"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- 1. DEFINISIKAN MOCKS ---

// mockUserRepository
type mockUserRepository struct {
	ports.UserRepository
	MockFindByEmail func(ctx context.Context, email string) (*domain.User, error)
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return m.MockFindByEmail(ctx, email)
}

// mockOtpRepository
type mockOtpRepository struct {
	ports.OTPRepository
	MockSave func(ctx context.Context, otp *domain.OTP) error
}

func (m *mockOtpRepository) Save(ctx context.Context, otp *domain.OTP) error {
	return m.MockSave(ctx, otp)
}

// mockEmailService
type mockEmailService struct {
	ports.EmailService
	MockSendPasswordResetEmail func(ctx context.Context, userEmail, name, otp string) error
}

func (m *mockEmailService) SendPasswordResetEmail(ctx context.Context, userEmail, name, otp string) error {
	return m.MockSendPasswordResetEmail(ctx, userEmail, name, otp)
}

// --- 2. FUNGSI SETUP ---

// setup_auth_service diperbarui
func setup_auth_service(_ *testing.T) (ports.AuthService, *mockUserRepository, *mockOtpRepository, *mockEmailService) {
	userRepo := &mockUserRepository{}
	otpRepo := &mockOtpRepository{}
	emailSvc := &mockEmailService{}
	dummyTokenCfg := config.Token{}
	nopLogger := zerolog.Nop()

	// Panggil konstruktor menggunakan prefix package: 'service.NewAuthService'
	service := service.NewAuthService(
		userRepo,
		nil,
		nil,
		dummyTokenCfg,
		nopLogger,
		nil,
		otpRepo,
		emailSvc,
	)

	// Kita tidak perlu konversi 'ok', karena NewAuthService sudah
	// mengembalikan interface ports.AuthService
	return service, userRepo, otpRepo, emailSvc
}

// --- 3. TES KASUS ---

// Skenario 1: Happy Path
func TestForgotPassword_HappyPath(t *testing.T) {
	// --- Arrange (Persiapan) ---
	// Tipe 'sut' sekarang adalah interface ports.AuthService, bukan *authService
	sut, mockUserRepo, mockOtpRepo, mockEmailSvc := setup_auth_service(t)

	testEmail := "user@example.com"
	testUser := &domain.User{
		ID:    1,
		UUID:  "user-uuid-123",
		Name:  "Test User",
		Email: testEmail,
	}

	var capturedOTP *domain.OTP
	var capturedEmailOTP string
	var capturedEmailTo string

	mockUserRepo.MockFindByEmail = func(ctx context.Context, email string) (*domain.User, error) {
		assert.Equal(t, testEmail, email)
		return testUser, nil
	}
	mockOtpRepo.MockSave = func(ctx context.Context, otp *domain.OTP) error {
		capturedOTP = otp
		return nil
	}
	mockEmailSvc.MockSendPasswordResetEmail = func(ctx context.Context, userEmail, name, otp string) error {
		capturedEmailTo = userEmail
		capturedEmailOTP = otp
		return nil
	}

	// --- Act (Tindakan) ---
	err := sut.ForgotPassword(context.Background(), testEmail)

	// --- Assert (Penegasan) ---
	require.NoError(t, err)
	require.NotNil(t, capturedOTP, "otpRepo.Save() seharusnya dipanggil")
	require.NotEmpty(t, capturedEmailOTP, "emailSvc.SendPasswordResetEmail() seharusnya dipanggil")
	assert.Equal(t, testUser.Email, capturedEmailTo)
	assert.Equal(t, testUser.ID, *capturedOTP.UserID)
	assert.Equal(t, domain.OTPTypePasswordReset, capturedOTP.Type)
	assert.Len(t, capturedOTP.Code, 6)
	assert.Equal(t, capturedOTP.Code, capturedEmailOTP)
	assert.WithinDuration(t, time.Now().Add(15*time.Minute), capturedOTP.ExpiresAt, 1*time.Minute)
}

// Skenario 2: User Tidak Ditemukan
func TestForgotPassword_UserNotFound(t *testing.T) {
	sut, mockUserRepo, mockOtpRepo, mockEmailSvc := setup_auth_service(t)
	testEmail := "nobody@example.com"
	otpSaveCalled := false
	emailSendCalled := false

	mockUserRepo.MockFindByEmail = func(ctx context.Context, email string) (*domain.User, error) {
		return nil, ports.ErrUserNotFound
	}
	mockOtpRepo.MockSave = func(ctx context.Context, otp *domain.OTP) error {
		otpSaveCalled = true
		return nil
	}
	mockEmailSvc.MockSendPasswordResetEmail = func(ctx context.Context, userEmail, name, otp string) error {
		emailSendCalled = true
		return nil
	}

	err := sut.ForgotPassword(context.Background(), testEmail)

	require.NoError(t, err)
	assert.False(t, otpSaveCalled)
	assert.False(t, emailSendCalled)
}

// Skenario 3: Gagal Menyimpan OTP
func TestForgotPassword_OtpSaveFails(t *testing.T) {
	sut, mockUserRepo, mockOtpRepo, mockEmailSvc := setup_auth_service(t)
	testUser := &domain.User{ID: 1, Email: "user@example.com", Name: "Test User"}
	dbError := errors.New("database connection error")
	emailSendCalled := false

	mockUserRepo.MockFindByEmail = func(ctx context.Context, email string) (*domain.User, error) {
		return testUser, nil
	}
	mockOtpRepo.MockSave = func(ctx context.Context, otp *domain.OTP) error {
		return dbError
	}
	mockEmailSvc.MockSendPasswordResetEmail = func(ctx context.Context, userEmail, name, otp string) error {
		emailSendCalled = true
		return nil
	}

	err := sut.ForgotPassword(context.Background(), "user@example.com")

	require.Error(t, err)
	assert.Equal(t, dbError, err)
	assert.False(t, emailSendCalled)
}

// Skenario 4: Gagal Mengirim Email
func TestForgotPassword_EmailSendFails(t *testing.T) {
	sut, mockUserRepo, mockOtpRepo, mockEmailSvc := setup_auth_service(t)
	testUser := &domain.User{ID: 1, Email: "user@example.com", Name: "Test User"}
	emailError := errors.New("mailgun API down")
	otpSaveCalled := false

	mockUserRepo.MockFindByEmail = func(ctx context.Context, email string) (*domain.User, error) {
		return testUser, nil
	}
	mockOtpRepo.MockSave = func(ctx context.Context, otp *domain.OTP) error {
		otpSaveCalled = true
		return nil
	}
	mockEmailSvc.MockSendPasswordResetEmail = func(ctx context.Context, userEmail, name, otp string) error {
		return emailError
	}

	err := sut.ForgotPassword(context.Background(), "user@example.com")

	require.NoError(t, err)
	assert.True(t, otpSaveCalled)
}
