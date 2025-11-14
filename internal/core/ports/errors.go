package ports

import "errors"

// Port-level errors define the contract for known business and data logic failures.
// This decouples the core service from any specific adapter (e.g., database) implementation.
var (
	// --- Auth & User Errors ---
	ErrEmailExists          = errors.New("email already in use")
	ErrInvalidCredentials   = errors.New("invalid email or password")
	ErrUserNotFound         = errors.New("user not found")
	ErrAccountDisabled      = errors.New("account is disabled")
	ErrEmailNotVerified     = errors.New("email is not verified")
	ErrEmailAlreadyVerified = errors.New("email is already verified")

	// --- Token Errors ---
	ErrInvalidToken     = errors.New("invalid token")
	ErrTokenExpired     = errors.New("token has expired")
	ErrTokenUseMismatch = errors.New("token cannot be used for this purpose")
	ErrTokenNotFound    = errors.New("token not found")

	// --- RBAC Errors ---
	ErrRoleNotFound       = errors.New("role not found")
	ErrPermissionNotFound = errors.New("permission not found")

	// --- General Errors ---
	ErrValidationFailed = errors.New("validation failed")

	// --- Repository Errors ---
	// These are translated by the adapter from driver-specific errors.
	ErrDuplicateEmail = errors.New("duplicate key (email) violates unique constraint")
	ErrUnexpectedSave = errors.New("unexpected error during user save")
	ErrUnexpectedFind = errors.New("unexpected error during user find")
)
