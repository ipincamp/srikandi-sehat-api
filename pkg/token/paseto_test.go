package token_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ipincamp/srikandi-sehat/pkg/token"
)

// Constants for testing
const (
	testSymmetricKey = "12345678901234567890123456789012" // 32 bytes
	testIssuer       = "srikandi-sehat-test"
	testUserID       = "user-uuid-12345"
	testRoleID       = "role-uuid-67890"
)

// TestPasetoMaker_Success_AccessToken tests the happy path for an Access Token.
// 1. Create a token.
// 2. Validate the token.
// 3. Ensure payloads match.
func TestPasetoMaker_Success_AccessToken(t *testing.T) {
	// The test must now pass the issuer, as the constructor changed.
	maker, err := token.NewPasetoMaker(testSymmetricKey, testIssuer)
	require.NoError(t, err)

	useFor := token.UseForAccessToken // Use the package constant
	duration := 15 * time.Minute

	// 1. Create Token
	tokenString, payload, err := maker.CreateToken(testUserID, testRoleID, useFor, duration)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)
	require.NotNil(t, payload)

	// Check payload from creation
	assert.Equal(t, testUserID, payload.UserID)
	assert.Equal(t, testRoleID, payload.RoleID)
	assert.Equal(t, useFor, payload.UseFor)
	assert.WithinDuration(t, time.Now(), payload.IssuedAt, time.Second)
	assert.WithinDuration(t, time.Now().Add(duration), payload.ExpiresAt, time.Second)

	// 2. Validate Token
	validatedPayload, err := maker.ValidateToken(tokenString)
	require.NoError(t, err)
	require.NotNil(t, validatedPayload)

	// Check payload from validation
	assert.Equal(t, payload.UserID, validatedPayload.UserID)
	assert.Equal(t, payload.RoleID, validatedPayload.RoleID)
	assert.Equal(t, payload.UseFor, validatedPayload.UseFor)
	assert.Equal(t, payload.ExpiresAt.Unix(), validatedPayload.ExpiresAt.Unix())
	assert.Equal(t, payload.IssuedAt.Unix(), validatedPayload.IssuedAt.Unix())
}

// TestPasetoMaker_Success_RefreshToken tests the happy path for a Refresh Token.
// This demonstrates that the RoleID can be empty.
func TestPasetoMaker_Success_RefreshToken(t *testing.T) {
	maker, err := token.NewPasetoMaker(testSymmetricKey, testIssuer)
	require.NoError(t, err)

	useFor := token.UseForRefreshToken // Use the constant
	roleID := ""                       // Refresh tokens don't need a role
	duration := 7 * 24 * time.Hour     // 7 days

	// 1. Create Token
	tokenString, payload, err := maker.CreateToken(testUserID, roleID, useFor, duration)
	require.NoError(t, err)
	require.NotNil(t, payload)

	// Check payload
	assert.Equal(t, testUserID, payload.UserID)
	assert.Equal(t, "", payload.RoleID)
	assert.Equal(t, useFor, payload.UseFor)

	// 2. Validate
	validatedPayload, err := maker.ValidateToken(tokenString)
	require.NoError(t, err)
	assert.Equal(t, testUserID, validatedPayload.UserID)
	assert.Equal(t, useFor, validatedPayload.UseFor)
	assert.Equal(t, "", validatedPayload.RoleID)
}

// TestPasetoMaker_ExpiredToken tests that validating an expired token
// correctly returns a `token.ErrTokenExpired`.
func TestPasetoMaker_ExpiredToken(t *testing.T) {
	maker, err := token.NewPasetoMaker(testSymmetricKey, testIssuer)
	require.NoError(t, err)

	// Create a token with a negative duration (i.e., expired 5 minutes ago)
	tokenString, _, err := maker.CreateToken(testUserID, testRoleID, token.UseForAccessToken, -5*time.Minute)
	require.NoError(t, err)

	// Try to validate it
	payload, err := maker.ValidateToken(tokenString)
	require.Error(t, err)
	assert.Nil(t, payload)

	// This is the most important check:
	// We must get our package's public error, not a raw Paseto error.
	assert.ErrorIs(t, err, token.ErrTokenExpired, "Expected token expired error")
}

// TestNewPasetoMaker_InvalidKeySize tests the constructor's validation.
func TestNewPasetoMaker_InvalidKeySize(t *testing.T) {
	// Key is too short
	shortKey := "tooshort"
	maker, err := token.NewPasetoMaker(shortKey, testIssuer)
	require.Error(t, err)
	assert.Nil(t, maker)
	assert.ErrorIs(t, err, token.ErrInvalidKeySize, "Expected invalid key size error for short key")

	// Key is too long
	longKey := "thiskeyiswaytoolongandshouldbeexactly32byteslong"
	maker, err = token.NewPasetoMaker(longKey, testIssuer)
	require.Error(t, err)
	assert.Nil(t, maker)
	assert.ErrorIs(t, err, token.ErrInvalidKeySize, "Expected invalid key size error for long key")
}

// TestPasetoMaker_InvalidTokenString tests validation against a
// string that is not a valid Paseto token.
func TestPasetoMaker_InvalidTokenString(t *testing.T) {
	maker, err := token.NewPasetoMaker(testSymmetricKey, testIssuer)
	require.NoError(t, err)

	// A plausible-looking but completely invalid token string
	invalidToken := "v2.local.this-is-not-a-valid-paseto-token-at-all.c29tZXNhbHQ"

	payload, err := maker.ValidateToken(invalidToken)
	require.Error(t, err)
	assert.Nil(t, payload)
	// Should be a parsing/decryption error, *not* an expired error
	assert.NotErrorIs(t, err, token.ErrTokenExpired, "Should be a parsing error, not expired error")
}
