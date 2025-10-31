package factories

import (
	"database/sql"
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

	return models.User{
		Name:  factory.Name,
		Email: factory.Email,
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

	hashedPassword, err := utils.HashPassword("password")
	if err != nil {
		return nil, err
	}

	var authProviders []models.UserAuthProvider
	for _, user := range users {
		authProviders = append(authProviders, models.UserAuthProvider{
			UserID:   user.ID,
			Provider: "local",
			Password: sql.NullString{String: hashedPassword, Valid: true},
		})
	}

	if err := db.CreateInBatches(&authProviders, 100).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func CreateAdminUser(tx *gorm.DB) (models.User, error) {
	// 1. Buat/Cari User
	adminUser := models.User{
		Name:  config.Get("ADMIN_NAME"),
		Email: config.Get("ADMIN_EMAIL"),
	}
	if err := tx.Where(models.User{Email: adminUser.Email}).FirstOrCreate(&adminUser).Error; err != nil {
		return models.User{}, err
	}

	// 2. Buat/Cari Auth Provider untuk admin
	adminPassword, err := utils.HashPassword(config.Get("ADMIN_PASSWORD"))
	if err != nil {
		return models.User{}, err
	}

	authProvider := models.UserAuthProvider{
		UserID:   adminUser.ID,
		Provider: "local",
		Password: sql.NullString{String: adminPassword, Valid: true},
	}
	// Cari berdasarkan UserID dan Provider
	if err := tx.Where(models.UserAuthProvider{UserID: adminUser.ID, Provider: "local"}).
		Attrs(authProvider). // Atribut ini dipakai jika FirstOrCreate membuat record baru
		FirstOrCreate(&authProvider).Error; err != nil {
		return models.User{}, err
	}

	return adminUser, nil
}
