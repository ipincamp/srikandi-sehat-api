package workers

import (
	"errors"
	"fmt"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/models/menstrual"
	"ipincamp/srikandi-sehat/src/utils"
	"time"

	"gorm.io/gorm"
)

const latePeriodThresholdDays = 32

// CheckLateMenstrualCycles finds users whose next cycle is late and sends a notification.
func CheckLateMenstrualCycles() {
	utils.InfoLogger.Println("Running Job: CheckLateMenstrualCycles...")
	var users []models.User
	// Ambil semua user yang memiliki FCM token
	database.DB.Where("fcm_token IS NOT NULL AND fcm_token != ?", "").Find(&users)

	for _, user := range users {
		var latestCycle menstrual.MenstrualCycle
		// Dapatkan siklus terakhir dari pengguna
		err := database.DB.Where("user_id = ?", user.ID).Order("start_date desc").First(&latestCycle).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue // Pengguna belum punya data siklus, lewati.
			}
			utils.ErrorLogger.Printf("Error fetching latest cycle for user %d: %v\n", user.ID, err)
			continue
		}

		// Proses hanya jika siklus terakhir sudah selesai (tidak sedang haid)
		// dan notifikasi keterlambatan belum pernah dikirim untuk siklus ini.
		if latestCycle.EndDate.Valid && !latestCycle.LatePeriodNotified {
			durationSinceEnd := time.Since(latestCycle.EndDate.Time)
			daysSinceEnd := int(durationSinceEnd.Hours() / 24)

			if daysSinceEnd > latePeriodThresholdDays {
				title := "Peringatan Keterlambatan Siklus"
				body := fmt.Sprintf("Sudah %d hari sejak siklus terakhir Anda selesai dan siklus baru belum dimulai. Segera periksakan diri jika Anda khawatir.", daysSinceEnd)

				err := utils.SendFCMNotification(user.ID, user.FcmToken, title, body, nil)
				if err != nil {
					utils.ErrorLogger.Printf("Failed to send late cycle notification to user %d: %v\n", user.ID, err)
					continue
				}

				// Tandai agar tidak mengirim notifikasi berulang untuk keterlambatan yang sama
				database.DB.Model(&latestCycle).Update("late_period_notified", true)
				utils.InfoLogger.Printf("Sent late cycle notification to user %d.", user.ID)
			}
		}
	}
	utils.InfoLogger.Println("Job: CheckLateMenstrualCycles finished.")
}
