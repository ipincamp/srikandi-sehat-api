package domain

// Permissions is the core domain model for a permission.
// This struct represents the business entity, free of any
// database or transport layer details.
type Permission struct {

	// ID is the internal database identifier.
	ID string

	// Name is the unique name of the permission.
	Name string

	// Description provides additional details about the permission.
	Description string
}
