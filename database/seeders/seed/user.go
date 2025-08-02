package seed

import (
	"errors"
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
			Name:     "Akun Administrator",
			Email:    config.Get("ADMIN_EMAIL"),
			Password: config.Get("USER_PASSWORD"),
			RoleName: constants.AdminRole,
		},
		{
			Name:     "Akun Pengguna",
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

		var existingUser models.User
		err := tx.Where("email = ?", userData.Email).First(&existingUser).Error

		if err == nil {
			continue
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
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

		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		if err := tx.Model(&user).Association("Roles").Replace(&role); err != nil {
			return err
		}
	}

	log.Println("User seeder completed successfully.")
	return nil
}
