package factories

import (
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/utils"

	"github.com/go-faker/faker/v4"
	"gorm.io/gorm"
)

func MakeUser() (models.User, error) {
	var factory models.User
	if err := faker.FakeData(&factory); err != nil {
		return models.User{}, err
	}

	hashedPassword, err := utils.HashPassword("password")
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		Name:     factory.Name,
		Email:    factory.Email,
		Password: hashedPassword,
	}, nil
}

func CreateUser(db *gorm.DB) (models.User, error) {
	user, err := MakeUser()
	if err != nil {
		return models.User{}, err
	}

	if err := db.Create(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func CreateUsers(db *gorm.DB, count int) ([]models.User, error) {
	var users []models.User
	for range count {
		user, err := MakeUser()
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := db.Create(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
