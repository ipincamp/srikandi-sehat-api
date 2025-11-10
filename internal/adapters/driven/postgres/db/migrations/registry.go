package migrations

import "github.com/go-gormigrate/gormigrate/v2"

// AllMigrations is the list of all migrations in the system.
// WHEN YOU CREATE A NEW MIGRATION, ADD IT TO THIS LIST.
// The 'make create-migration' command will create the file,
// but you must manually add the constructor function call here.
var AllMigrations = []*gormigrate.Migration{
	// Add new migrations here, e.g.:
	CreateUsersTable(),
	CreateOtpsTable(),
	// CreateRolesTable(),
	// CreateProductsTable(),
}
