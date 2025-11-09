package password

// Hasher is an interface for hashing and comparing passwords.
// By defining an interface, we apply the Dependency Inversion Principle (DIP).
// Our business logic (e.g., auth use cases) will depend on this
// abstraction, not on a concrete implementation like Argon2 or Bcrypt.
// This makes our system pluggable and easy to migrate to a new
// hashing algorithm in the future without changing business logic.
type Hasher interface {
	// Hash generates a hashed version of the given password.
	Hash(password string) (string, error)

	// Compare checks if the given password matches the hashed version.
	// It's crucial that this function *does not* return an error.
	// Any processing error (e.g., invalid hash format) should
	// simply result in `false` (no match).
	Compare(hash string, password string) bool
}
