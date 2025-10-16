package workers

import (
	"fmt"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/models/menstrual"
	"ipincamp/srikandi-sehat/src/utils"
	"time"
)

const longPeriodThresholdDays = 7

// CheckLongMenstrualCycles finds ongoing cycles longer than the threshold and sends a notification.
func CheckLongMenstrualCycles() {
	utils.InfoLogger.Println("Running Job: CheckLongMenstrualCycles...")
	var activeCycles []menstrual.MenstrualCycle

	// 1. Ambil semua siklus yang masih aktif (end_date is null) & belum pernah dinotifikasi
	err := database.DB.Preload("User").
		Where("menstrual_cycles.end_date IS NULL AND menstrual_cycles.long_period_notified = ?", false).
		Find(&activeCycles).Error

	if err != nil {
		utils.ErrorLogger.Printf("Error fetching active cycles for checker: %v\n", err)
		return
	}

	for _, cycle := range activeCycles {
		// 2. Hitung durasi dari tanggal mulai hingga hari ini
		duration := time.Since(cycle.StartDate)
		days := int(duration.Hours() / 24)

		// 3. Jika durasi > 7 hari, kirim notifikasi
		if days > longPeriodThresholdDays {
			// Pastikan user punya FCM token
			if cycle.User.FcmToken == "" {
				utils.InfoLogger.Printf("User %d has no FCM token, skipping long cycle notification.", cycle.UserID)
				continue
			}

			// Siapkan dan kirim notifikasi
			title := "Peringatan Durasi Menstruasi"
			body := fmt.Sprintf("Siklus menstruasi Anda saat ini sudah berlangsung selama %d hari. Batas normalnya adalah 3-7 hari.", days)
			err := utils.SendFCMNotification(cycle.UserID, cycle.User.FcmToken, title, body, nil)

			if err != nil {
				utils.ErrorLogger.Printf("Failed to send long cycle notification to user %d: %v\n", cycle.UserID, err)
				continue // Lanjut ke user berikutnya meskipun gagal
			}

			// 4. Update flag di database agar tidak dikirim notifikasi lagi
			database.DB.Model(&cycle).Update("long_period_notified", true)
			utils.InfoLogger.Printf("Sent long cycle notification to user %d for cycle ID %d.", cycle.UserID, cycle.ID)
		}
	}
	utils.InfoLogger.Println("Job: CheckLongMenstrualCycles finished.")
}
