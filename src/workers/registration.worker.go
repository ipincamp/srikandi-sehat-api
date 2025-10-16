package workers

import (
	"errors"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/constants"
	"ipincamp/srikandi-sehat/src/dto"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/utils"
	"log"
	"runtime"

	"gorm.io/gorm"
)

type Job struct {
	RegistrationData dto.RegisterRequest
	FCMToken         string
}

var JobQueue chan Job

func StartWorkerPool() {
	numWorkers := runtime.NumCPU()
	queueSize := 5000
	JobQueue = make(chan Job, queueSize)

	for w := 1; w <= numWorkers; w++ {
		go worker(w, JobQueue)
	}

	log.Printf("Starting %d workers for registration process...", numWorkers)
}

func worker(id int, jobs <-chan Job) {
	for job := range jobs {
		data := job.RegistrationData
		fcmToken := job.FCMToken

		var existingUser models.User
		err := database.DB.First(&existingUser, "email = ?", data.Email).Error
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Worker %d: Registration failed for %s. Email is already registered.", id, data.Email)
			utils.SendFCMNotification(
				0,
				fcmToken,
				"Registration Failed",
				"The email you used is already registered in our system.",
				map[string]string{"status": "failed", "reason": "email_exists"},
			)
			continue
		}

		hashedPassword, err := utils.HashPassword(data.Password)
		if err != nil {
			log.Printf("ERROR (Worker %d): Failed to hash password for %s: %v", id, data.Email, err)
			continue
		}

		user := models.User{
			Name:     data.Name,
			Email:    data.Email,
			Password: hashedPassword,
		}

		err = database.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&user).Error; err != nil {
				return err
			}

			defaultRole, err := utils.GetRoleByName(string(constants.UserRole))
			if err != nil {
				return err
			}

			if err := tx.Model(&user).Association("Roles").Append(&defaultRole); err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			log.Printf("ERROR (Worker %d): Failed to create user %s in database: %v", id, data.Email, err)
			utils.SendFCMNotification(
				user.ID,
				fcmToken,
				"Registration Failed",
				"An error occurred while processing your account. Please try again.",
				map[string]string{"status": "failed", "reason": "server_error"},
			)
		} else {
			utils.AddEmailToRegistrationFilter(user.Email)
			log.Printf("Worker %d: User %s has been created and activated.", id, user.Email)
			utils.SendFCMNotification(
				user.ID,
				fcmToken,
				"Registration Successful!",
				"Your account has been created successfully. Please log in to get started.",
				map[string]string{"status": "success"},
			)
		}
	}
}
