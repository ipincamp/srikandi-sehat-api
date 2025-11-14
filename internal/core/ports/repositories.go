package ports

import (
	"context"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
)

// UserRepository defines the contract for user-related data operations.
// This repository handles all user data persistence including CRUD operations,
// role assignments, and permission management. It serves as the driven port
// for user data access in the hexagonal architecture.
type UserRepository interface {
	// Save persists a new user or updates an existing user in the repository.
	// This method should handle both insert and update operations based on whether
	// the user has an existing ID. Returns an error if the operation fails or if
	// there's a constraint violation (e.g., duplicate email).
	Save(ctx context.Context, user *domain.User) error

	// FindByID retrieves a user by their unique UUID.
	// Returns ErrUserNotFound if no user exists with the given UUID.
	// This method is commonly used for profile retrieval and authentication checks.
	FindByID(ctx context.Context, uuid string) (*domain.User, error)

	// FindByEmail retrieves a user by their email address.
	// Returns ErrUserNotFound if no user exists with the given email.
	// This method is primarily used during login and registration validation.
	FindByEmail(ctx context.Context, email string) (*domain.User, error)

	// FindMapByUUIDs returns a map of users keyed by UUID for the given list of UUIDs.
	// This method is optimized for batch operations and is typically used by DataLoader
	// to prevent N+1 query problems. Missing UUIDs are simply omitted from the result map.
	FindMapByUUIDs(ctx context.Context, uuids []string) (map[string]*domain.User, error)

	// Update modifies the details of an existing user.
	// Only the provided fields in the user object will be updated.
	// Returns ErrUserNotFound if the user doesn't exist.
	Update(ctx context.Context, user *domain.User) error

	// Delete removes a user from the repository by their UUID.
	// This may perform a soft delete (marking as deleted) or hard delete depending
	// on implementation. Returns ErrUserNotFound if the user doesn't exist.
	Delete(ctx context.Context, uuid string) error

	// AssignRoleToUser assigns a role to a user by their respective IDs.
	// This creates a many-to-many relationship between users and roles.
	// Returns an error if either the user or role doesn't exist, or if the
	// assignment already exists.
	AssignRoleToUser(ctx context.Context, userID string, roleID string) error

	// RemoveRoleFromUser removes a specific role from a user.
	// This deletes the many-to-many relationship between the user and role.
	// Returns an error if the user or role doesn't exist, or if the assignment
	// doesn't exist.
	RemoveRoleFromUser(ctx context.Context, userID string, roleID string) error

	// ListUserRoles lists all roles assigned to a user.
	// Returns an empty slice if the user has no roles.
	// This is used for authorization checks and user management interfaces.
	ListUserRoles(ctx context.Context, userID string) ([]*domain.Role, error)

	// AddDirectPermissionToUser grants a direct permission to a user.
	// Direct permissions are assigned to the user independently of their roles.
	// This is useful for granting exceptional permissions outside of role hierarchies.
	// Returns an error if the user or permission doesn't exist, or if the permission
	// is already directly assigned.
	AddDirectPermissionToUser(ctx context.Context, userID string, permissionID string) error

	// RemoveDirectPermissionFromUser revokes a direct permission from a user.
	// This only removes direct permissions, not those inherited through roles.
	// Returns an error if the user or permission doesn't exist, or if the direct
	// permission assignment doesn't exist.
	RemoveDirectPermissionFromUser(ctx context.Context, userID string, permissionID string) error

	// ListUserDirectPermissions lists all direct permissions assigned to a user.
	// This does NOT include permissions inherited through roles.
	// Returns an empty slice if the user has no direct permissions.
	ListUserDirectPermissions(ctx context.Context, userID string) ([]*domain.Permission, error)

	// ListUserAllPermissions lists all permissions (direct and via roles) assigned to a user.
	// This aggregates both direct permissions and permissions inherited from all assigned roles,
	// removing duplicates. This method is the primary way to check what a user is authorized to do.
	// Returns an empty slice if the user has no permissions.
	ListUserAllPermissions(ctx context.Context, userID string) ([]*domain.Permission, error)
}

// UserTokenRepository manages temporary tokens for various user operations.
// This repository handles OTP codes, password reset tokens, email verification tokens,
// and other time-limited, single-use tokens. Tokens are typically consumed after use.
type UserTokenRepository interface {
	// Save persists a new user token to the repository.
	// Each token should have a unique hash, purpose (e.g., "password_reset", "email_verification"),
	// and expiration time. Returns an error if the save operation fails.
	Save(ctx context.Context, token *domain.UserToken) error

	// FindAndConsume retrieves a token matching the given criteria and marks it as consumed.
	// This is an atomic operation that ensures the token can only be used once.
	// Returns ErrTokenNotFound if no matching token exists, ErrTokenExpired if the token
	// has passed its expiration time, or ErrTokenAlreadyUsed if already consumed.
	FindAndConsume(ctx context.Context, userID string, purpose string, tokenHash string) (*domain.UserToken, error)
}

// RoleRepository manages role definitions and their associated permissions.
// Roles are collections of permissions that can be assigned to users, implementing
// role-based access control (RBAC). Common examples include "admin", "user", "moderator".
type RoleRepository interface {
	// Save persists a new role or updates an existing role in the repository.
	// Role names should be unique. Returns an error if the save fails or if
	// there's a duplicate name constraint violation.
	Save(ctx context.Context, role *domain.Role) error

	// FindByID retrieves a role by its unique identifier.
	// Returns ErrRoleNotFound if no role exists with the given ID.
	FindByID(ctx context.Context, roleID string) (*domain.Role, error)

	// FindByName retrieves a role by its name.
	// Role names are typically unique identifiers like "admin" or "moderator".
	// Returns ErrRoleNotFound if no role exists with the given name.
	FindByName(ctx context.Context, name string) (*domain.Role, error)

	// Delete removes a role from the repository by its ID.
	// This should also remove all user-role and role-permission associations.
	// Returns ErrRoleNotFound if the role doesn't exist.
	Delete(ctx context.Context, roleID string) error

	// List retrieves all roles in the system.
	// Returns an empty slice if no roles exist.
	// This is typically used for administration interfaces.
	List(ctx context.Context) ([]*domain.Role, error)

	// AddPermissionToRole associates a permission with a role.
	// All users with this role will inherit the permission.
	// Returns an error if the role or permission doesn't exist, or if the
	// association already exists.
	AddPermissionToRole(ctx context.Context, roleID string, permissionID string) error

	// RemovePermissionFromRole removes a permission from a role.
	// This affects all users with this role.
	// Returns an error if the role or permission doesn't exist, or if the
	// association doesn't exist.
	RemovePermissionFromRole(ctx context.Context, roleID string, permissionID string) error
}

// PermissionRepository manages permission definitions in the system.
// Permissions represent specific actions or access rights, such as "create_post",
// "delete_user", or "view_analytics". They are the atomic units of authorization.
type PermissionRepository interface {
	// FindByID retrieves a permission by its unique identifier.
	// Returns ErrPermissionNotFound if no permission exists with the given ID.
	FindByID(ctx context.Context, permissionID string) (*domain.Permission, error)

	// FindByName retrieves a permission by its name.
	// Permission names should be unique and descriptive, like "users.create" or "posts.delete".
	// Returns ErrPermissionNotFound if no permission exists with the given name.
	FindByName(ctx context.Context, name string) (*domain.Permission, error)

	// List retrieves all permissions defined in the system.
	// Returns an empty slice if no permissions exist.
	// This is typically used for administration interfaces and permission management.
	List(ctx context.Context) ([]*domain.Permission, error)

	// Save persists a new permission or updates an existing permission.
	// Permission names should be unique. Returns an error if the save fails or if
	// there's a duplicate name constraint violation.
	Save(ctx context.Context, permission *domain.Permission) error
}

// RefreshTokenRepository manages long-lived API tokens for programmatic access.
// Personal access tokens (PATs) allow users to authenticate API requests without
// using their password, similar to GitHub's personal access tokens. These are typically
// used for CLI tools, scripts, and third-party integrations.
type RefreshTokenRepository interface {
	// Save persists a new personal access token to the repository.
	// The token should be hashed before storage for security.
	// Returns an error if the save operation fails.
	Save(ctx context.Context, token *domain.PersonalToken) error

	// FindByID retrieves a personal access token by its unique identifier.
	// This is typically called after matching a hashed token to verify its metadata.
	// Returns ErrTokenNotFound if no token exists with the given ID.
	FindByID(ctx context.Context, tokenID string) (*domain.PersonalToken, error)

	// ListByUserID retrieves all personal access tokens belonging to a user.
	// Returns an empty slice if the user has no tokens.
	// This is used to display a user's active tokens in their account settings.
	ListByUserID(ctx context.Context, userID string) ([]*domain.PersonalToken, error)

	// Delete removes a personal access token from the repository.
	// The userID parameter ensures users can only delete their own tokens.
	// Returns ErrTokenNotFound if the token doesn't exist or doesn't belong to the user.
	Delete(ctx context.Context, tokenID string, userID string) error
}

// ActivityLogRepository manages audit logs for user activities.
// This repository stores records of important user actions for security auditing,
// compliance, and debugging purposes. Typical logged activities include logins,
// password changes, permission modifications, and sensitive data access.
type ActivityLogRepository interface {
	// Save persists a new activity log entry to the repository.
	// Activity logs should be append-only and never modified or deleted for audit integrity.
	// Returns an error if the save operation fails.
	Save(ctx context.Context, log *domain.ActivityLog) error

	// ListByUserID retrieves all activity logs for a specific user.
	// Logs are typically returned in reverse chronological order (newest first).
	// Returns an empty slice if no logs exist for the user.
	// Consider implementing pagination for users with many log entries.
	ListByUserID(ctx context.Context, userID string) ([]*domain.ActivityLog, error)
}
