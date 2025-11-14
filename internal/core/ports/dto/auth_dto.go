package dto

// -- Request DTOs ---

// LoginAuthRequest represents the data required for a user login operation.
type LoginAuthRequest struct {
	// Email is the user's email address.
	Email string

	// Password is the user's password.
	Password string
}

// RefreshTokenAuthRequest represents the data required to refresh authentication tokens.
type RefreshTokenAuthRequest struct {
	// RefreshToken is the token used to obtain new authentication tokens.
	RefreshToken string
}

// LogoutAuthRequest represents the data required to logout a user.
type LogoutAuthRequest struct {
	// RefreshToken is the token to be invalidated during logout.
	RefreshToken string
}

// -- Response DTOs ---

// TokenAuthResponse represents the response containing access and refresh tokens.
type TokenAuthResponse struct {
	// AccessToken is access token.
	AccessToken string

	// RefreshToken is refresh token.
	RefreshToken string
}
