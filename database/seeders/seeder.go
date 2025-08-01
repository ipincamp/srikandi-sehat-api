package main

import "gorm.io/gorm"

func SeedAll(db *gorm.DB) error {
	if err := SeedRoles(db); err != nil {
		return err
	}
	if err := SeedPermissions(db); err != nil {
		return err
	}
	if err := SeedUsers(db); err != nil {
		return err
	}

	return nil
}
