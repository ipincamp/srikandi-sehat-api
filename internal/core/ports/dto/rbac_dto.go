package dto

// -- Request DTOs ---

// CreateRoleUserRequest represents the data required to create a new role.
type CreateRoleUserRequest struct {
	// Name is the name of the role to be created.
	Name string

	// Description is an optional description of the role.
	Description *string
}

// AssignPermissionToRoleUserRequest represents the data required to add permissions to a role.
type AssignPermissionToRoleUserRequest struct {
	// RoleID is the unique identifier of the role.
	RoleID string

	// PermissionIDs is a list of permission IDs to be associated with the role.
	PermissionIDs []string
}

// RemovePermissionFromRoleUserRequest represents the data required to remove permissions from a role.
type RemovePermissionFromRoleUserRequest struct {
	// RoleID is the unique identifier of the role.
	RoleID string

	// PermissionIDs is a list of permission IDs to be disassociated from the role.
	PermissionIDs []string
}

// AssignRoleUserRequest represents the data required to assign a role to a user.
type AssignRoleUserRequest struct {
	// UserID is the unique identifier of the user.
	IdUserRequest

	// RoleID is the unique identifier of the role to be assigned.
	RoleID string
}

// RemoveRoleUserRequest represents the data required to remove a role from a user.
type RemoveRoleUserRequest struct {
	// UserID is the unique identifier of the user.
	IdUserRequest

	// RoleID is the unique identifier of the role to be removed.
	RoleID string
}

// AssignDirectPermissionUserRequest represents the data required to assign direct permissions to a user.
type AssignDirectPermissionUserRequest struct {
	// UserID is the unique identifier of the user.
	IdUserRequest

	// PermissionIDs is a list of permission IDs to be directly assigned to the user.
	PermissionIDs []string
}

// RemoveDirectPermissionUserRequest represents the data required to remove direct permissions from a user.
type RemoveDirectPermissionUserRequest struct {
	// UserID is the unique identifier of the user.
	IdUserRequest

	// PermissionIDs is a list of permission IDs to be removed from the user.
	PermissionIDs []string
}
