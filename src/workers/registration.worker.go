package workers

import (
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/constants"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/utils"
	"log"
	"runtime"

	"gorm.io/gorm"
)

type Job struct {
	User          models.User
	PlainPassword string
}

var JobQueue chan Job

func StartWorkerPool() {
	numWorkers := runtime.NumCPU()
	JobQueue = make(chan Job, 1000)

	for w := 1; w <= numWorkers; w++ {
		go worker(w, JobQueue)
	}

	log.Printf("Starting %d workers for registration process...", numWorkers)
}

func worker(id int, jobs <-chan Job) {
	for job := range jobs {
		log.Printf("Worker %d: starting process for user %s", id, job.User.Email)

		hashedPassword, err := utils.HashPassword(job.PlainPassword)
		if err != nil {
			log.Printf("ERROR (Worker %d): Failed to hash password for user %s: %v", id, job.User.Email, err)
			continue
		}

		var defaultRole models.Role
		if err := database.DB.First(&defaultRole, "name = ?", string(constants.UserRole)).Error; err != nil {
			log.Printf("CRITICAL (Worker %d): Default role '%s' not found for user %s", id, constants.UserRole, job.User.Email)
			continue
		}

		err = database.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&job.User).Updates(models.User{Password: hashedPassword, Status: constants.StatusActive}).Error; err != nil {
				return err
			}
			if err := tx.Model(&job.User).Association("Roles").Append(&defaultRole); err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			log.Printf("ERROR (Worker %d): Failed to update user %s after hashing: %v", id, job.User.Email, err)
		} else {
			utils.AddEmailToFilter(job.User.Email)
			log.Printf("Worker %d: Process for user %s completed.", id, job.User.Email)
		}
	}
}
