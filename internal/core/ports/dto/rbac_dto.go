package dto

// -- Request DTOs ---

// CreateRoleRequest represents the data required to create a new role.
type CreateRoleRequest struct {
	// Name is the name of the role to be created.
	Name string

	// Description is an optional description of the role.
	Description *string
}

type CreatePermissionRequest struct {
	// Name is the name of the permission to be created.
	Name string

	// Description is an optional description of the permission.
	Description *string
}

// IdRoleRequest represents a request containing only a role ID.
type IdRoleRequest struct {
	// ID is the unique identifier of the role.
	ID string
}

// AssignPermissionsToRoleRequest represents the data required to assign permissions to a role.
type AssignPermissionsToRoleRequest struct {
	// RoleID is the unique identifier of the role.
	IdRoleRequest

	// PermissionIDs is a list of permission IDs to be associated with the role.
	PermissionIDs []string
}

// RemovePermissionsFromRoleRequest represents the data required to remove permissions from a role.
type RemovePermissionsFromRoleRequest struct {
	// RoleID is the unique identifier of the role.
	IdRoleRequest

	// PermissionIDs is a list of permission IDs to be disassociated from the role.
	PermissionIDs []string
}

// AssignRolesToUserRequest represents the data required to assign roles to a user.
type AssignRolesToUserRequest struct {
	// UserID is the unique identifier of the user.
	IdUserRequest

	// RoleIDs is a list of role IDs to be assigned to the user.
	RoleIDs []string
}

// RemoveRolesFromUserRequest represents the data required to remove roles from a user.
type RemoveRolesFromUserRequest struct {
	// UserID is the unique identifier of the user.
	IdUserRequest

	// RoleIDs is a list of role IDs to be removed from the user.
	RoleIDs []string
}

// AssignDirectPermissionsToUserRequest represents the data required to assign direct permissions to a user.
type AssignDirectPermissionsToUserRequest struct {
	// UserID is the unique identifier of the user.
	IdUserRequest

	// PermissionIDs is a list of permission IDs to be directly assigned to the user.
	PermissionIDs []string
}

// RemoveDirectPermissionsFromUserRequest represents the data required to remove direct permissions from a user.
type RemoveDirectPermissionsFromUserRequest struct {
	// UserID is the unique identifier of the user.
	IdUserRequest

	// PermissionIDs is a list of permission IDs to be removed from the user.
	PermissionIDs []string
}
