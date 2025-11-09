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

// Make sure Argon2idHasher implements the Hasher interface
var _ Hasher = (*Argon2idHasher)(nil)

// Argon2idHasher is a Hasher implementation struct.
// It uses the Argon2id algorithm for hashing passwords.
type Argon2idHasher struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// NewArgon2idHasher creates a new instance of Argon2idHasher
// with safe default parameters.
func NewArgon2idHasher() Hasher {
	return &Argon2idHasher{
		Memory:      64 * 1024, // 64 MB
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
}

// Hash generates an Argon2id hash of the password.
func (h *Argon2idHasher) Hash(password string) (string, error) {
	// 1. Generate cryptographically secure salt
	salt := make([]byte, h.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// 2. Generate hash
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		h.Iterations,
		h.Memory,
		h.Parallelism,
		h.KeyLength,
	)

	// 3. Encode to standard string format for storing in DB
	// Format: $argon2id$v=19$m=<memory>,t=<iterations>,p=<parallelism>$<salt>$<hash>

	// Base64 encode salt and hash
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Format string
	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		h.Memory,
		h.Iterations,
		h.Parallelism,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

// Compare verifies the plaintext password against the Argon2id hash.
func (h *Argon2idHasher) Compare(encodedHash string, password string) bool {
	// 1. Parse the hash string to get the parameters, salt, and hash.
	params, salt, hash, err := h.decodeHash(encodedHash)
	if err != nil {
		// If the hash format is wrong, just assume it doesn't match.
		// Never panic or return error here.
		return false
	}

	// 2. Generate a *new* hash from the given password,
	// using the *exact same parameters and salt* from the old hash.
	otherHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		// KeyLength is obtained from the length of the decoded hash.
		params.KeyLength,
	)

	// 3. Compare the two hashes using constant-time compare
	// to prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true
	}
	return false
}

// -- Helper Internal --

type argonParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	// We get this from the actual hash length
	KeyLength uint32
}

// decodeHash parses the string hash format
func (h *Argon2idHasher) decodeHash(encodedHash string) (p *argonParams, salt, hash []byte, err error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, errors.New("invalid hash format: wrong number of parts")
	}

	if parts[1] != "argon2id" {
		return nil, nil, nil, errors.New("invalid hash format: not argon2id")
	}

	var version int
	_, err = fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil || version != argon2.Version {
		return nil, nil, nil, errors.New("invalid hash format: incompatible version")
	}

	p = &argonParams{}
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism)
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

	p.KeyLength = uint32(len(hash))

	return p, salt, hash, nil
}
