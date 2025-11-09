package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Make sure Argon2idHasher implements the Hasher interface.
// This is a compile-time check that ensures our struct
// correctly satisfies the contract.
var _ Hasher = (*Argon2idHasher)(nil)

// Argon2idHasher is a Hasher implementation struct.
// It uses the Argon2id algorithm.
//
// The fields are *unexported* (lowercase) to enforce encapsulation.
// This prevents other packages from creating an instance with insecure
// parameters (e.g., &Argon2idHasher{memory: 1024}).
// All instances *must* be created via the NewArgon2idHasher() constructor.
type Argon2idHasher struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// NewArgon2idHasher creates a new instance of Argon2idHasher
// with safe default parameters. These parameters are based on
// current (e.g., OWASP) recommendations.
// This constructor is the *only* way to create the struct,
// ensuring secure-by-default parameters.
func NewArgon2idHasher() Hasher {
	return &Argon2idHasher{
		memory:      64 * 1024, // 64 MB
		iterations:  3,
		parallelism: 2,
		saltLength:  16, // 16 bytes
		keyLength:   32, // 32 bytes
	}
}

// Hash generates an Argon2id hash of the password.
// The output is an MCF (Modular Crypt Format) compliant string.
func (h *Argon2idHasher) Hash(password string) (string, error) {
	// 1. Generate a cryptographically secure salt.
	salt := make([]byte, h.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// 2. Generate the hash using the parameters from the struct.
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		h.iterations,
		h.memory,
		h.parallelism,
		h.keyLength,
	)

	// 3. Encode to standard string format for storing in the DB.
	// Format: $argon2id$v=19$m=<memory>,t=<iterations>,p=<parallelism>$<salt>$<hash>

	// Base64 encode salt and hash using Raw Standard Encoding
	// (no padding, as specified by the MCF).
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// 4. Format the final string.
	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		h.memory,
		h.iterations,
		h.parallelism,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

// Compare verifies the plaintext password against the Argon2id hash.
// This function is robust against malformed hashes and safe
// against timing attacks.
func (h *Argon2idHasher) Compare(encodedHash string, password string) bool {
	// 1. Parse the hash string to get the parameters, salt, and hash.
	params, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		// If the hash format is wrong, it's not a match.
		// Never panic or return an error from Compare.
		return false
	}

	// 2. Generate a *new* hash from the *given* password,
	// using the *exact same parameters and salt* from the *old* hash.
	// This is the only correct way to verify a password.
	otherHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.iterations,
		params.memory,
		params.parallelism,
		params.keyLength, // KeyLength is derived from the *decoded hash length*
	)

	// 3. Compare the two hashes using a constant-time comparison
	// to prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true
	}
	return false
}

// -- Helper Internal --

// argonParams holds the decoded parameters from a hash string.
// This is *unexported* (lowercase 'a') as it's an internal
// implementation detail of this package.
type argonParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	keyLength   uint32
}

// Added constants to remove magic numbers/strings from the decodeHash function.
const (
	expectedHashParts   = 6
	algorithmIdentifier = "argon2id"
)

// decodeHash parses the string hash format.
// This is *unexported* (lowercase 'd') as it's an internal
// implementation detail.
func decodeHash(encodedHash string) (p *argonParams, salt, hash []byte, err error) {
	// The hash string should be in the format:
	// $argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$c29tZXBhc3N3b3Jk
	//   [0]      [1]    [2]       [3]          [4]        [5]
	parts := strings.Split(encodedHash, "$")
	// Using constant instead of magic number 6.
	if len(parts) != expectedHashParts {
		return nil, nil, nil, errors.New("invalid hash format: wrong number of parts")
	}

	// Using constant instead of magic string "argon2id".
	if parts[1] != algorithmIdentifier {
		return nil, nil, nil, errors.New("invalid hash format: not argon2id")
	}

	var version int
	_, err = fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil || version != argon2.Version {
		return nil, nil, nil, errors.New("invalid hash format: incompatible version")
	}

	p = &argonParams{}
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, errors.New("invalid hash format: bad parameters")
	}

	salt, err = base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, errors.New("invalid hash format: bad salt")
	}

	hash, err = base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, errors.New("invalid hash format: bad hash")
	}

	// The key length is derived from the *actual* length of the
	// decoded hash, not from the struct's default.
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}
