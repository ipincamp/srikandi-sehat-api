package utils

import (
	"context"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/models"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

var fcmClient *messaging.Client

func InitFCM() {
	opt := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		ErrorLogger.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		ErrorLogger.Fatalf("error getting Messaging client: %v\n", err)
	}

	fcmClient = client
	InfoLogger.Println("Firebase Cloud Messaging client initialized successfully.")
}

func SendFCMNotification(userID uint, token string, title string, body string, data map[string]string) error {
	// --- Tahap 1: Coba kirim Push Notification (FCM) ---
	if fcmClient == nil {
		ErrorLogger.Println("FCM client is not initialized. Skipping FCM send, proceeding to save history.")
	} else if token != "" {
		// Hanya coba kirim jika token ada
		message := &messaging.Message{
			Notification: &messaging.Notification{
				Title: title,
				Body:  body,
			},
			Data:  data,
			Token: token,
		}

		_, err := fcmClient.Send(context.Background(), message)
		if err != nil {
			// Jika FCM gagal: Catat error, tapi JANGAN return.
			// Tetap menyimpan notifikasi ke DB.
			ErrorLogger.Printf("Failed to send FCM notification to token %s: %v. Proceeding to save history.", token, err)
		} else {
			// Jika FCM berhasil: Catat sukses.
			InfoLogger.Printf("Successfully sent FCM notification to token %s", token)
		}
	} else {
		// Jika user tidak punya token: Catat info, lanjut simpan ke DB.
		InfoLogger.Printf("Skipping FCM send for user %d (no token), saving to DB history only.", userID)
	}

	// --- Tahap 2: Simpan notifikasi ke database (History) ---
	// Tahap ini sekarang SELALU dijalankan, baik FCM sukses maupun gagal.
	notification := models.Notification{
		UserID: userID,
		Title:  title,
		Body:   body,
	}
	if err := database.DB.Create(&notification).Error; err != nil {
		// Jika DB GAGAL: Ini adalah error serius. Catat dan kembalikan error.
		ErrorLogger.Printf("Failed to save notification to database for user %d: %v", userID, err)
		return err // Mengembalikan error DB
	}

	// Jika sampai di sini, artinya history DB berhasil disimpan.
	// Operasi dianggap sukses dari perspektif worker.
	return nil
}
