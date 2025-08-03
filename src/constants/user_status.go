package constants

type UserStatus string

const (
	StatusProcessing UserStatus = "processing"
	StatusActive     UserStatus = "active"
	StatusSuspended  UserStatus = "suspended"
)
