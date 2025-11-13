package domain

// Role is the core domain model for a user role.
// This struct represents the business entity, free of any
// database or transport layer details.
type Role struct {

	// ID is the internal database identifier.
	ID string

	// Name is the unique name of the role.
	Name string

	// Description provides additional details about the role.
	Description string

	// Permissions is the list of permissions associated with the role.
	Permissions []*Permission
}
