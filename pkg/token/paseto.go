package token

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/o1egl/paseto"
)

// This compile-time check ensures that our *pasetoMaker struct
// always satisfies the Maker interface. If we break the interface
// (e.g., change a function signature), the compiler will fail,
// protecting us from runtime errors.
var _ Maker = (*pasetoMaker)(nil)

// --- Package-Level Errors ---

// ErrInvalidKeySize is returned from the constructor if the key
// is not the required 32-byte length for Paseto's V2 symmetric mode.
var ErrInvalidKeySize = errors.New("invalid symmetric key size: must be exactly 32 bytes")

// ErrInvalidToken is returned by ValidateToken if the token is malformed,
// fails decryption, or is missing essential custom claims.
var ErrInvalidToken = errors.New("token is invalid or missing claims")

// (Note: ErrTokenExpired is defined in payload.go, as it's part of
// the public Payload contract).

// --- Paseto Maker Implementation ---

// pasetoMaker is the Paseto-specific implementation of the Maker interface.
// Its fields are unexported to enforce encapsulation.
type pasetoMaker struct {
	paseto       *paseto.V2 // The Paseto v2 library instance
	symmetricKey []byte     // The raw 32-byte symmetric key
	issuer       string     // The name of this service (e.g., "srikandi-sehat-api")
}

// NewPasetoMaker creates a new instance of pasetoMaker.
// This constructor acts as a "guard" to ensure the maker is
// created with valid parameters.
func NewPasetoMaker(symmetricKey string, issuer string) (Maker, error) {
	// 1. Fail-fast: Validate the key size immediately.
	//    This prevents a runtime panic inside the paseto library.
	if len(symmetricKey) != 32 {
		return nil, ErrInvalidKeySize
	}

	// 2. Return the fully initialized struct, satisfying the Maker interface.
	return &pasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
		issuer:       issuer,
	}, nil
}

// CreateToken generates a new PASETO v2.local token.
func (m *pasetoMaker) CreateToken(userID, roleID, useFor string, duration time.Duration) (string, *Payload, error) {
	// 1. Create the custom payload struct. This validates inputs (like UserID)
	//    and sets the 'iat' and 'exp' times.
	payload, err := NewPayload(userID, roleID, useFor, duration)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create payload: %w", err)
	}

	// 2. Create the Paseto standard JSON token, filling in standard claims.
	jsonToken := paseto.JSONToken{
		Issuer:     m.issuer,          // 'iss' claim
		IssuedAt:   payload.IssuedAt,  // 'iat' claim
		Expiration: payload.ExpiresAt, // 'exp' claim
		Subject:    payload.UserID,    // 'sub' claim
	}

	// 3. Add our application-specific custom claims.
	//    These will be stored in the token's JSON payload.
	jsonToken.Set("uid", payload.UserID)
	jsonToken.Set("rid", payload.RoleID)
	jsonToken.Set("for", payload.UseFor)

	// 4. Encrypt the token using Paseto v2.local (symmetric key).
	token, err := m.paseto.Encrypt(m.symmetricKey, jsonToken, nil)
	if err != nil {
		return "", nil, fmt.Errorf("failed to encrypt token: %w", err)
	}

	return token, payload, nil
}

// ValidateToken decrypts and verifies a PASETO v2.local token.
func (m *pasetoMaker) ValidateToken(token string) (*Payload, error) {
	// 1. We decrypt into a standard paseto.JSONToken.
	var jsonToken paseto.JSONToken
	var footer interface{} // We are not using a footer, so this is nil.

	// 2. Decrypt the token. The Paseto library automatically verifies the
	//    token's structure and *some* standard claims (like 'exp'
	//    if it's in the past, though we double-check later).
	err := m.paseto.Decrypt(token, m.symmetricKey, &jsonToken, &footer)
	if err != nil {
		// Error Translation:
		// If Paseto returns a "token has expired" error, we translate
		// it into our package's public ErrTokenExpired.
		// This decouples the caller from the Paseto library.
		if strings.Contains(err.Error(), "token has expired") {
			return nil, ErrTokenExpired
		}
		// For all other decryption errors (bad format, bad signature),
		// return our public ErrInvalidToken.
		return nil, fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}

	// 3. Manually extract our custom claims from the token's internal map.
	//    We use .Get() which returns a string.
	userID := jsonToken.Get("uid")
	roleID := jsonToken.Get("rid")
	useFor := jsonToken.Get("for")

	// 4. Sanity check our *required* custom claims.
	//    A valid token *must* contain a user ID and its intended use.
	if userID == "" || useFor == "" {
		return nil, ErrInvalidToken
	}

	// 5. Re-create our public Payload struct from the claims.
	payload := &Payload{
		UserID:    userID,
		RoleID:    roleID,
		UseFor:    useFor,
		IssuedAt:  jsonToken.IssuedAt,
		ExpiresAt: jsonToken.Expiration,
	}

	// 6. ***CRITICAL***: Perform our own explicit validation.
	//    This is a defense-in-depth step. It re-checks the expiration
	//    using time.Now(), fixing the 'TestPasetoMaker_ExpiredToken' failure.
	if err := payload.Valid(); err != nil {
		return nil, err // This will correctly return ErrTokenExpired
	}

	return payload, nil
}
