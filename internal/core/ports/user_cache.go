package ports

// UserCache defines the "driven port" for an in-memory user cache.
// This is used by core services to quickly check for the *possible*
// existence of a user before hitting the main database.
type UserCache interface {
	// Add adds a user's email to the cache.
	// This is typically called after successful registration.
	Add(email string)

	// Test checks if a user's email *might* exist.
	// - Returns false if the email is *definitely not* in the cache.
	// - Returns true if the email is *possibly* in the cache (false positive).
	Test(email string) bool
}
