package migrations

import "github.com/go-gormigrate/gormigrate/v2"

// AllMigrations is the list of all migrations in the system.
// WHEN YOU CREATE A NEW MIGRATION, ADD IT TO THIS LIST.
var AllMigrations = []*gormigrate.Migration{
	// Add new migrations here, e.g.:
	CreateUsersTable(),
	// CreateRolesTable(),
	// CreateProductsTable(),
}
