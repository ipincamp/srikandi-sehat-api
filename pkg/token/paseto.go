package token

import (
	"errors"
	"fmt"
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

// --- Added Constants ---

// pasetoV2SymmetricKeySize is the required 32-byte key length for PASETO v2.local.
const pasetoV2SymmetricKeySize = 32

// Custom claim keys. Using constants prevents typos and ensures
// CreateToken and ValidateToken use the same values.
const (
	claimUserID = "uid"
	claimRoleID = "rid"
	claimUseFor = "for"
)

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
	//    Using constant instead of magic number 32.
	if len(symmetricKey) != pasetoV2SymmetricKeySize {
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
		Jti:        payload.JTI,       // 'jti' claim
	}

	// 3. Add our application-specific custom claims.
	//    These will be stored in the token's JSON payload.
	//    Using constants instead of magic strings.
	jsonToken.Set(claimUserID, payload.UserID)
	jsonToken.Set(claimRoleID, payload.RoleID)
	jsonToken.Set(claimUseFor, payload.UseFor)

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
		//  REMOVED fragile string-based error check.
		// The `payload.Validate()` check at the end of this function
		// is the robust way to check for expiration.
		//
		// Old code:
		// if strings.Contains(err.Error(), "token has expired") {
		// 	return nil, ErrTokenExpired
		// }

		// For all other decryption errors (bad format, bad signature),
		// return our public ErrInvalidToken.
		return nil, fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}

	// 3. Manually extract our custom claims from the token's internal map.
	//    We use .Get() which returns a string.
	//    Using constants instead of magic strings.
	userID := jsonToken.Get(claimUserID)
	roleID := jsonToken.Get(claimRoleID)
	useFor := jsonToken.Get(claimUseFor)

	// 4. Sanity check our *required* custom claims.
	//    A valid token *must* contain a user ID and its intended use.
	if userID == "" || useFor == "" {
		return nil, ErrInvalidToken
	}

	// 5. Re-create our public Payload struct from the claims.
	payload := &Payload{
		JTI:       jsonToken.Jti,
		UserID:    userID,
		RoleID:    roleID,
		UseFor:    useFor,
		IssuedAt:  jsonToken.IssuedAt,
		ExpiresAt: jsonToken.Expiration,
	}

	// 6. ***CRITICAL***: Perform our own explicit validation.
	//    This is a defense-in-depth step. It re-checks the expiration
	//    using time.Now(), fixing the 'TestPasetoMaker_ExpiredToken' failure.
	//    Using the renamed `Validate()` method.
	if err := payload.Validate(); err != nil {
		return nil, err // This will correctly return ErrTokenExpired
	}

	return payload, nil
}
