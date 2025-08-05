package factories

import (
	"ipincamp/srikandi-sehat/config"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/utils"

	"github.com/go-faker/faker/v4"
	"gorm.io/gorm"
)

type userFactory struct {
	Name  string `faker:"name"`
	Email string `faker:"email"`
}

func MakeUser() (models.User, error) {
	var factory userFactory
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

func CreateUsers(db *gorm.DB, count int) ([]models.User, error) {
	var users []models.User
	for range count {
		user, err := MakeUser()
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := db.CreateInBatches(&users, 100).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func CreateAdminUser(tx *gorm.DB) (models.User, error) {
	adminPassword, err := utils.HashPassword(config.Get("ADMIN_PASSWORD"))
	if err != nil {
		return models.User{}, err
	}

	adminUser := models.User{
		Name:     config.Get("ADMIN_NAME"),
		Email:    config.Get("ADMIN_EMAIL"),
		Password: adminPassword,
	}

	if err := tx.Where(models.User{Email: adminUser.Email}).FirstOrCreate(&adminUser).Error; err != nil {
		return models.User{}, err
	}

	return adminUser, nil
}
