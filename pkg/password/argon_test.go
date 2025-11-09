package password_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ipincamp/srikandi-sehat/pkg/password"
)

// TestHashAndCompare_Success tests the happy path:
// 1. Hash a password.
// 2. Compare the *same* password against the hash.
// 3. Expect a match.
func TestHashAndCompare_Success(t *testing.T) {
	hasher := password.NewArgon2idHasher()
	pass := "myS3cureP@ssword123!"

	// 1. Create hash
	hashedPass, err := hasher.Hash(pass)

	// 'require' will stop the test on failure (Fatal)
	require.NoError(t, err, "Hashing should not produce an error")
	require.NotEmpty(t, hashedPass, "Hash should not be empty")

	// 2. Ensure hash is not the same as the password
	require.NotEqual(t, pass, hashedPass, "Hash must not be equal to the password")

	// 3. Compare
	// 'assert' will record a failure but continue the test
	match := hasher.Compare(hashedPass, pass)
	assert.True(t, match, "The correct password should match")
}

// TestCompare_Failure_WrongPassword tests that a *wrong* password
// does not match a valid hash.
func TestCompare_Failure_WrongPassword(t *testing.T) {
	hasher := password.NewArgon2idHasher()
	pass := "myS3cureP@ssword123!"
	wrongPass := "thisIsTheWrongPassword"

	hashedPass, err := hasher.Hash(pass)
	require.NoError(t, err)

	// Compare with the wrong password
	match := hasher.Compare(hashedPass, wrongPass)
	assert.False(t, match, "The wrong password should not match")
}

// TestCompare_Failure_InvalidHash tests that the Compare function
// correctly returns `false` for various malformed hash strings.
func TestCompare_Failure_InvalidHash(t *testing.T) {
	hasher := password.NewArgon2idHasher()
	pass := "password123"

	// Test various invalid hash formats
	invalidHashes := []string{
		"just-random-string",                                             // Random text
		"$argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ",                     // Missing parts
		"bad$argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$c29tZXBhc3N3b3Jk", // Wrong prefix
		"$argon2id$v=18$m=65536,t=3,p=2$c29tZXNhbHQ$c29tZXBhc3N3b3Jk",    // Wrong version
		"$argon2id$v=19$m=xx,t=3,p=2$c29tZXNhbHQ$c29tZXBhc3N3b3Jk",       // Corrupt parameters
		"$argon2id$v=19$m=65536,t=3,p=2$!!!$c29tZXBhc3N3b3Jk",            // Corrupt base64 salt
		"$argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$!!!",                 // Corrupt base64 hash
	}

	for _, invalidHash := range invalidHashes {
		// Use a subtest for better failure reporting
		t.Run(invalidHash, func(t *testing.T) {
			match := hasher.Compare(invalidHash, pass)
			assert.False(t, match, "Invalid hash format '%s' should not match", invalidHash)
		})
	}
}

// TestHash_CreatesDifferentHashes tests that the salting mechanism works.
// Hashing the *same* password twice should produce two *different*
// hash strings.
func TestHash_CreatesDifferentHashes(t *testing.T) {
	hasher := password.NewArgon2idHasher()
	pass := "same-password-over-and-over"

	// Generate two hashes for the same password
	hash1, err1 := hasher.Hash(pass)
	hash2, err2 := hasher.Hash(pass)

	require.NoError(t, err1)
	require.NoError(t, err2)

	// Because the salt is random, the resulting hashes must be different
	assert.NotEqual(t, hash1, hash2, "Two hashes for the same password must be different due to salt")

	// However, both hashes must still be valid for the original password
	assert.True(t, hasher.Compare(hash1, pass), "Hash 1 must be valid")
	assert.True(t, hasher.Compare(hash2, pass), "Hash 2 must be valid")
}
