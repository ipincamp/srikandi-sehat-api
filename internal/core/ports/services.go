package ports

import (
	"context"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
	"github.com/ipincamp/srikandi-sehat/internal/core/ports/dto"
)

// =============================================================================
// Authentication Service
// =============================================================================

// AuthService defines the contract for authentication operations including
// identity verification and token management.
type AuthService interface {
	// Login validates user credentials and returns authentication tokens on success.
	Login(ctx context.Context, dto dto.LoginAuthRequest) (*dto.TokenAuthResponse, error)

	// RefreshToken validates a refresh token and issues a new token pair.
	RefreshToken(ctx context.Context, dto dto.RefreshTokenAuthRequest) (*dto.TokenAuthResponse, error)

	// Logout invalidates a refresh token.
	Logout(ctx context.Context, dto dto.LogoutAuthRequest) error
}

// =============================================================================
// User Service
// =============================================================================

// UserService defines the contract for comprehensive user management including
// lifecycle operations, profile management, and RBAC.
type UserService interface {
	// Register creates a new user account and returns authentication tokens.
	Register(ctx context.Context, dto dto.RegisterUserRequest) (*dto.TokenAuthResponse, error)

	// ForgotPassword initiates password reset flow by generating and sending an OTP to the user's email.
	ForgotPassword(ctx context.Context, dto dto.ForgotPasswordUserRequest) error

	// ResetPassword completes password reset flow by validating OTP and updating user's password.
	ResetPassword(ctx context.Context, dto dto.ResetPasswordUserRequest) error

	// --- User Lifecycle ---

	// GetUserByID retrieves a user by their unique UUID.
	GetUserByID(ctx context.Context, dto dto.IdUserRequest) (*domain.User, error)

	// ListAllUsers retrieves all users in the system with pagination.
	ListAllUsers(ctx context.Context, dto dto.UserFilterRequest) ([]*domain.User, int, error)

	// DisableAccount marks an account as disabled, preventing login while retaining data for potential reactivation.
	DisableAccount(ctx context.Context, dto dto.IdUserRequest) error

	// DeleteAccount permanently removes a user account (may perform soft delete based on implementation).
	DeleteAccount(ctx context.Context, dto dto.IdUserRequest) error

	// --- Profile & Security ---

	// UpdateProfile updates user profile information such as name.
	UpdateProfile(ctx context.Context, dto dto.UpdateProfileUserRequest) (*domain.User, error)

	// ChangePassword updates the password for an authenticated user after verifying the old password.
	ChangePassword(ctx context.Context, dto dto.ChangePasswordUserRequest) error

	// --- Email Verification ---

	// RequestEmailVerification initiates email verification flow by generating and sending an OTP to the user.
	RequestEmailVerification(ctx context.Context, dto dto.IdUserRequest) error

	// VerifyEmail completes email verification flow by validating OTP and marking email as verified.
	VerifyEmail(ctx context.Context, dto dto.VerifyEmailUserRequest) error

	// --- RBAC (Role-Based Access Control) ---

	// AssignPermissionsToRole assigns permissions to a role.
	AssignPermissionsToRole(ctx context.Context, dto dto.AssignPermissionsToRoleRequest) error

	// RemovePermissionsFromRole revokes permissions from a role.
	RemovePermissionsFromRole(ctx context.Context, dto dto.RemovePermissionsFromRoleRequest) error

	// AssignRolesToUser assigns roles to a user.
	AssignRolesToUser(ctx context.Context, dto dto.AssignRolesToUserRequest) error

	// RemoveRolesFromUser revokes roles from a user.
	RemoveRolesFromUser(ctx context.Context, dto dto.RemoveRolesFromUserRequest) error

	// AssignDirectPermissionsToUser grants direct permissions to a user.
	AssignDirectPermissionsToUser(ctx context.Context, dto dto.AssignDirectPermissionsToUserRequest) error

	// RemoveDirectPermissionsFromUser revokes direct permissions from a user.
	RemoveDirectPermissionsFromUser(ctx context.Context, dto dto.RemoveDirectPermissionsFromUserRequest) error
}

type RBACService interface {

	// CreateRole creates a new role with the specified name and description.
	CreateRole(ctx context.Context, dto dto.CreateRoleRequest) (*domain.Role, error)

	// ListAllRoles retrieves all roles in the system.
	ListAllRoles(ctx context.Context, dto dto.RoleFilterRequest) ([]*domain.Role, error)

	// CreatePermission creates a new permission with the specified name and description.
	CreatePermission(ctx context.Context, dto dto.CreatePermissionRequest) (*domain.Permission, error)

	// ListAllPermissions retrieves all permissions defined in the system.
	ListAllPermissions(ctx context.Context, dto dto.PermissionFilterRequest) ([]*domain.Permission, error)
}

// =============================================================================
// Mail Service
// =============================================================================

// MailService defines the contract for email delivery adapter handling various
// email notifications. This is a driven port (outbound) for sending emails.
type MailService interface {
	// SendPasswordReset sends a password reset OTP email to the user.
	SendPasswordReset(ctx context.Context, dto dto.PasswordResetMailRequest) error

	// SendEmailVerification sends an email verification OTP to the user.
	SendEmailVerification(ctx context.Context, dto dto.EmailVerificationMailRequest) error

	// SendEmailChange sends an email change verification OTP to the new email address.
	SendEmailChange(ctx context.Context, dto dto.EmailChangeMailRequest) error

	// SendDeleteAccount sends an account deletion notification to the user.
	SendDeleteAccount(ctx context.Context, dto dto.AccountNotificationMailRequest) error

	// SendDisableAccount sends an account disabled notification to the user.
	SendDisableAccount(ctx context.Context, dto dto.AccountNotificationMailRequest) error
}
