package dto

// Request Body for testing FCM
type TestNotificationRequest struct {
	Title string            `json:"title" validate:"required,min=3"`
	Body  string            `json:"body" validate:"required,min=5"`
	Data  map[string]string `json:"data,omitempty"` // Optional data payload
}
