package domain

import "time"

// Action constants for logging
const (
	ActionLogin          = "LOGIN"
	ActionChangePassword = "CHANGE_PASSWORD"
	ActionUpdateProfile  = "UPDATE_PROFILE"
	// Add more actions as needed
)

// ActivityLog is the core domain model for an activity log entry.
// This struct represents the business entity, free of any
// database or transport layer details.
type ActivityLog struct {

	// ID is the internal database identifier.
	ID int64

	// ActorUserID is the ID of the user who performed the action.
	ActorUserID *string

	// Action describes the action taken.
	Action string

	// TargetTable is the name of the table affected by the action.
	TargetTable *string

	// TargetID is the ID of the record affected by the action.
	TargetID *string

	// Changes contains a JSON representation of the changes made.
	Changes *string

	// IPAddress is the IP address from which the action was performed.
	IPAddress *string

	// UserAgent is the user agent string of the client used to perform the action.
	UserAgent *string

	// RequestID is an optional identifier for tracing the request.
	RequestID *string

	// Timestamp is the time when the action was performed.
	Timestamp time.Time
}

// PaginatedActivityLogs is a response struct for the "GetMyLogs" use case.
type PaginatedActivityLogs struct {
	Logs       []*ActivityLog
	Pagination PaginationMetadata
}

// PaginationMetadata holds the metadata for a paginated response.
type PaginationMetadata struct {
	TotalItems  int64
	TotalPages  int
	CurrentPage int
	PerPage     int
}
