package domain

// AuthResponse is the core domain model for an authentication response.
// This contains the tokens returned to the user.
type AuthResponse struct {
	AccessToken  string
	RefreshToken string
}
