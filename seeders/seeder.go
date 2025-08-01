package main

import (
	"ipincamp/srikandi-sehat/config"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/constants"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/utils"
	"log"
)

func main() {
	config.LoadConfig()
	database.ConnectDB()

	log.Println("Seeding started...")

	seedRoles()
	seedPermissions()
	seedUsers()

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

func seedUsers() {
	type UserSeed struct {
		Name     string
		Email    string
		Password string
		RoleName constants.RoleName
	}

	usersToSeed := []UserSeed{
		{
			Name:     "Admin Name",
			Email:    config.Get("ADMIN_EMAIL"),
			Password: config.Get("USER_PASSWORD"),
			RoleName: constants.AdminRole,
		},
		{
			Name:     "User Name",
			Email:    config.Get("USER_EMAIL"),
			Password: config.Get("USER_PASSWORD"),
			RoleName: constants.UserRole,
		},
	}

	for _, userData := range usersToSeed {
		if userData.Email == "" || userData.Password == "" {
			log.Printf("Skipping seeding for role %s: email or password not set in .env", userData.RoleName)
			continue
		}

		hashedPassword, err := utils.HashPassword(userData.Password)
		if err != nil {
			log.Fatalf("Could not hash password for %s: %v", userData.Email, err)
		}

		var role models.Role
		if err := database.DB.First(&role, "name = ?", userData.RoleName).Error; err != nil {
			log.Fatalf("Role '%s' not found. Please run role seeder first.", userData.RoleName)
		}

		user := models.User{
			Name:     userData.Name,
			Email:    userData.Email,
			Password: hashedPassword,
		}

		var existingUser models.User
		if err := database.DB.FirstOrCreate(&existingUser, models.User{Email: user.Email}, &user).Error; err != nil {
			log.Fatalf("Could not seed user %s: %v", user.Email, err)
		}

		if err := database.DB.Model(&existingUser).Association("Roles").Replace(&role); err != nil {
			log.Fatalf("Could not assign role to user %s: %v", user.Email, err)
		}

		log.Printf("User '%s' with role '%s' seeded successfully.", existingUser.Email, role.Name)
	}
}
