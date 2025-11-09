package service

import "errors"

// Business logic errors
var (
	ErrEmailExists        = errors.New("email already in use")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token has expired")
	ErrTokenUseMismatch   = errors.New("token cannot be used for this purpose")
)
