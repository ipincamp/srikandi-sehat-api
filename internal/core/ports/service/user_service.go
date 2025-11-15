package service

import (
	"context"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
	"github.com/ipincamp/srikandi-sehat/internal/core/ports/dto"
)

// UserService defines the driving port for user management including CRUD operations, account status, and RBAC.
type UserService interface {
	// --- User Lifecycle ---

	// Register creates a new user account and returns authentication tokens.
	Register(ctx context.Context, dto dto.RegisterUserRequest) (*dto.TokenAuthResponse, error)

	// GetUserByID retrieves a user by their unique UUID.
	GetUserByID(ctx context.Context, dto dto.IdUserRequest) (*domain.User, error)

	// ForgotPassword initiates password reset flow by generating and sending an OTP to the user's email.
	ForgotPassword(ctx context.Context, dto dto.ForgotPasswordUserRequest) error

	// ResetPassword completes password reset flow by validating OTP and updating user's password.
	ResetPassword(ctx context.Context, dto dto.ResetPasswordUserRequest) error

	// DisableAccount marks an account as disabled, preventing login while retaining data for potential reactivation.
	DisableAccount(ctx context.Context, dto dto.IdUserRequest) error

	// DeleteAccount permanently removes a user account (may perform soft delete based on implementation).
	DeleteAccount(ctx context.Context, dto dto.IdUserRequest) error

	// ListAllUsers retrieves all users in the system with pagination.
	ListAllUsers(ctx context.Context, dto dto.UserFilterRequest) ([]*domain.User, int, error)

	// --- Profile & Verification ---

	// UpdateProfile updates user profile information such as name.
	UpdateProfile(ctx context.Context, dto dto.UpdateProfileUserRequest) (*domain.User, error)

	// ChangePassword updates the password for an authenticated user after verifying the old password.
	ChangePassword(ctx context.Context, dto dto.ChangePasswordUserRequest) error

	// RequestEmailVerification initiates email verification flow by generating and sending an OTP to the user.
	RequestEmailVerification(ctx context.Context, dto dto.IdUserRequest) error

	// VerifyEmail completes email verification flow by validating OTP and marking email as verified.
	VerifyEmail(ctx context.Context, dto dto.VerifyEmailUserRequest) error

	// --- RBAC (Role-Based Access Control) ---

	// CreateRole creates a new role with the specified name and description.
	CreateRole(ctx context.Context, dto dto.CreateRoleUserRequest) (*domain.Role, error)

	// ListAllRoles retrieves all roles in the system.
	ListAllRoles(ctx context.Context, dto dto.RoleFilterRequest) ([]*domain.Role, error)

	// AssignPermissionToRole associates a permission with a role.
	AssignPermissionToRole(ctx context.Context, dto dto.AssignPermissionToRoleUserRequest) error

	// RemovePermissionFromRole disassociates a permission from a role.
	RemovePermissionFromRole(ctx context.Context, dto dto.RemovePermissionFromRoleUserRequest) error

	// ListAllPermissions retrieves all permissions defined in the system.
	ListAllPermissions(ctx context.Context, dto dto.PermissionFilterRequest) ([]*domain.Permission, error)

	// AssignRoleToUser assigns a role to a user.
	AssignRoleToUser(ctx context.Context, dto dto.AssignRoleUserRequest) error

	// RemoveRoleFromUser removes a role from a user.
	RemoveRoleFromUser(ctx context.Context, dto dto.RemoveRoleUserRequest) error

	// AssignDirectPermissionToUser grants a direct permission to a user.
	AssignDirectPermissionToUser(ctx context.Context, dto dto.AssignDirectPermissionUserRequest) error

	// RemoveDirectPermissionFromUser revokes a direct permission from a user.
	RemoveDirectPermissionFromUser(ctx context.Context, dto dto.RemoveDirectPermissionUserRequest) error
}
