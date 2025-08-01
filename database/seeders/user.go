package main

import (
	"ipincamp/srikandi-sehat/config"
	"ipincamp/srikandi-sehat/src/constants"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/utils"
	"log"

	"gorm.io/gorm"
)

func SeedUsers(tx *gorm.DB) error {
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
			log.Printf("Skipping seeder for role %s: email or password not set in .env", userData.RoleName)
			continue
		}

		hashedPassword, err := utils.HashPassword(userData.Password)
		if err != nil {
			return err
		}

		var role models.Role
		if err := tx.First(&role, "name = ?", userData.RoleName).Error; err != nil {
			return err
		}

		user := models.User{
			Name:     userData.Name,
			Email:    userData.Email,
			Password: hashedPassword,
		}

		var existingUser models.User
		if err := tx.FirstOrCreate(&existingUser, models.User{Email: user.Email}, &user).Error; err != nil {
			return err
		}

		if err := tx.Model(&existingUser).Association("Roles").Replace(&role); err != nil {
			return err
		}

		log.Printf("User '%s' with role '%s' successfully seeded.", existingUser.Email, role.Name)
	}

	return nil
}
