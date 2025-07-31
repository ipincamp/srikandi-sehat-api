package main

import (
	"ipincamp/srikandi-sehat/config"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/constants"
	"ipincamp/srikandi-sehat/src/models"
	"log"
)

func main() {
	config.LoadConfig()
	database.ConnectDB()

	log.Println("Seeding started...")

	seedRoles()
	seedPermissions()

	log.Println("Seeding completed successfully!")
}

func seedRoles() {
	roles := []models.Role{
		{Name: string(constants.AdminRole)},
		{Name: string(constants.UserRole)},
	}

	for _, role := range roles {
		err := database.DB.FirstOrCreate(&role, models.Role{Name: role.Name}).Error
		if err != nil {
			log.Fatalf("Could not seed roles: %v", err)
		}
	}
	log.Println("Roles seeded")
}

func seedPermissions() {
	permissions := []models.Permission{
		// {Name: "view-profile"},
		// {Name: "edit-profile"},
	}

	for _, permission := range permissions {
		err := database.DB.FirstOrCreate(&permission, models.Permission{Name: permission.Name}).Error
		if err != nil {
			log.Fatalf("Could not seed permissions: %v", err)
		}
	}
	log.Println("Permissions seeded")
}
