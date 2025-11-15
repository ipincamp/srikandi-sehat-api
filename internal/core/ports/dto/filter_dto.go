package dto

// -- Request Query DTOs ---

// PaginateFilterRequest represents pagination parameters for listing requests.
type PaginateFilterRequest struct {
	// Page number to retrieve. Default is 1.
	Page *int

	// Number of items per page. Default is 10.
	Limit *int
}

// ListAllUserRequest represents the data required to list all users with pagination and filtering.
type UserFilterRequest struct {
	// Pagination parameters for listing users.
	PaginateFilterRequest

	// Name filters users by their name (partial match).
	Name *string

	// Email filters users by their email (partial match).
	Email *string

	// IsVerified filters users based on email verification status.
	IsVerified *bool

	// IsDisabled filters users based on account disabled status.
	IsDisabled *bool
}

// RoleFilterRequest represents the data required to list all roles with pagination and filtering.
type RoleFilterRequest struct {
	// Pagination parameters for listing roles.
	PaginateFilterRequest

	// Name filters roles by their name (partial match).
	Name *string

	// Description filters roles by their description (partial match).
	Description *string
}

// PermissionFilterRequest represents the data required to list all permissions with pagination and filtering.
type PermissionFilterRequest struct {
	// Pagination parameters for listing permissions.
	PaginateFilterRequest

	// Name filters permissions by their name (partial match).
	Name *string

	// Description filters permissions by their description (partial match).
	Description *string
}
