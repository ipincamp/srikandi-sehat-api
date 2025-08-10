package utils

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

var fcmClient *messaging.Client

func InitFCM() {
	opt := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	fcmClient = client
	InfoLogger.Println("Firebase Cloud Messaging client initialized successfully.")
}

func SendFCMNotification(token string, title string, body string, data map[string]string) error {
	if fcmClient == nil {
		ErrorLogger.Println("FCM client is not initialized.")
		return nil
	}

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
		ErrorLogger.Printf("Failed to send FCM notification to token %s: %v\n", token, err)
		return err
	}

	InfoLogger.Printf("Successfully sent FCM notification to token %s", token)
	return nil
}
