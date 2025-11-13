package ports

import (
	"context"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
)

// AuthResponse is a DTO containing access and refresh tokens returned by authentication operations.
type AuthResponse struct {
	AccessToken  string
	RefreshToken string
}

// AuthService defines the driving port for authentication operations including identity verification and token management.
type AuthService interface {
	// Register creates a new user account and returns authentication tokens.
	Register(ctx context.Context, name, email, password string) (*AuthResponse, error)

	// Login validates user credentials and returns authentication tokens on success.
	Login(ctx context.Context, email, password string) (*AuthResponse, error)

	// RefreshToken validates a refresh token and issues a new token pair.
	RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error)

	// Logout invalidates a refresh token (only relevant for stateful token implementations).
	Logout(ctx context.Context, refreshToken string) error

	// ForgotPassword initiates password reset flow by generating and sending an OTP to the user's email.
	ForgotPassword(ctx context.Context, email string) error

	// ResetPassword completes password reset flow by validating OTP and updating user's password.
	ResetPassword(ctx context.Context, email, otp, newPassword string) error

	// ChangePassword updates the password for an authenticated user after verifying the old password.
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error

	// RequestEmailVerification initiates email verification flow by generating and sending an OTP to the user.
	RequestEmailVerification(ctx context.Context, userID string) error

	// VerifyEmail completes email verification flow by validating OTP and marking email as verified.
	VerifyEmail(ctx context.Context, userID string, otp string) error
}

// UserService defines the driving port for user management including CRUD operations, account status, and RBAC.
type UserService interface {
	// GetUserByID retrieves a user by their unique UUID.
	GetUserByID(ctx context.Context, uuid string) (*domain.User, error)

	// UpdateProfile updates user profile information such as name.
	UpdateProfile(ctx context.Context, userID string, newName string) (*domain.User, error)

	// DisableAccount marks an account as disabled, preventing login while retaining data for potential reactivation.
	DisableAccount(ctx context.Context, userID string) error

	// DeleteAccount permanently removes a user account (may perform soft delete based on implementation).
	DeleteAccount(ctx context.Context, userID string) error

	// ListUsers retrieves all users in the system (admin only).
	ListUsers(ctx context.Context) ([]*domain.User, error)

	// CreateRole creates a new role with the specified name and description.
	CreateRole(ctx context.Context, name, description string) (*domain.Role, error)

	// ListRoles retrieves all roles in the system.
	ListRoles(ctx context.Context) ([]*domain.Role, error)

	// AddPermissionToRole associates a permission with a role.
	AddPermissionToRole(ctx context.Context, roleID, permissionID string) error

	// ListPermissions retrieves all permissions defined in the system.
	ListPermissions(ctx context.Context) ([]*domain.Permission, error)

	// AssignRoleToUser assigns a role to a user.
	AssignRoleToUser(ctx context.Context, userID, roleID string) error

	// RemoveRoleFromUser removes a role from a user.
	RemoveRoleFromUser(ctx context.Context, userID, roleID string) error

	// AddDirectPermissionToUser grants a direct permission to a user.
	AddDirectPermissionToUser(ctx context.Context, userID, permissionID string) error

	// RemoveDirectPermissionFromUser revokes a direct permission from a user.
	RemoveDirectPermissionFromUser(ctx context.Context, userID, permissionID string) error
}

// MailService defines the contract for email delivery adapter handling various email notifications.
type MailService interface {
	// SendPasswordResetEmail sends a password reset OTP email to the user.
	SendPasswordResetEmail(ctx context.Context, userEmail, name, otp string) error

	// SendEmailVerificationEmail sends an email verification OTP to the user.
	SendEmailVerificationEmail(ctx context.Context, userEmail, name, otp string) error

	// SendEmailChangeEmail sends an email change verification OTP to the new email address.
	SendEmailChangeEmail(ctx context.Context, oldEmail, newEmail, name, otp string) error

	// SendDeleteAccountEmail sends an account deletion notification to the user.
	SendDeleteAccountEmail(ctx context.Context, userEmail, name string) error

	// SendDisableAccountEmail sends an account disabled notification to the user.
	SendDisableAccountEmail(ctx context.Context, userEmail, name string) error
}
