package password_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ipincamp/srikandi-sehat/pkg/password"
)

// TestHashAndCompare_Success menguji happy path
func TestHashAndCompare_Success(t *testing.T) {
	hasher := password.NewArgon2idHasher()
	pass := "myS3cureP@ssword123!"

	// 1. Buat hash
	hashedPass, err := hasher.Hash(pass)

	// require akan menghentikan tes jika gagal (Fatal)
	require.NoError(t, err, "Hashing seharusnya tidak error")
	require.NotEmpty(t, hashedPass, "Hash tidak boleh kosong")

	// 2. Pastikan hash tidak sama dengan password
	require.NotEqual(t, pass, hashedPass, "Hash tidak boleh sama dengan password")

	// 3. Bandingkan
	// assert hanya akan mencatat kegagalan tapi melanjutkan tes
	match := hasher.Compare(hashedPass, pass)
	assert.True(t, match, "Password yang benar seharusnya cocok")
}

// TestCompare_Failure_WrongPassword menguji password yang salah
func TestCompare_Failure_WrongPassword(t *testing.T) {
	hasher := password.NewArgon2idHasher()
	pass := "myS3cureP@ssword123!"
	wrongPass := "iniPasswordYangSalah"

	hashedPass, err := hasher.Hash(pass)
	require.NoError(t, err)

	// Bandingkan dengan password yang salah
	match := hasher.Compare(hashedPass, wrongPass)
	assert.False(t, match, "Password yang salah seharusnya tidak cocok")
}

// TestCompare_Failure_InvalidHash menguji format hash yang rusak
func TestCompare_Failure_InvalidHash(t *testing.T) {
	hasher := password.NewArgon2idHasher()
	pass := "password123"

	// Tes berbagai format hash yang salah
	invalidHashes := []string{
		"just-random-string",                                             // Teks acak
		"$argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ",                     // Kurang bagian
		"bad$argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$c29tZXBhc3N3b3Jk", // Awalan salah
		"$argon2id$v=18$m=65536,t=3,p=2$c29tZXNhbHQ$c29tZXBhc3N3b3Jk",    // Versi salah
		"$argon2id$v=19$m=xx,t=3,p=2$c29tZXNhbHQ$c29tZXBhc3N3b3Jk",       // Parameter rusak
		"$argon2id$v=19$m=65536,t=3,p=2$!!!$c29tZXBhc3N3b3Jk",            // Salt base64 rusak
		"$argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$!!!",                 // Hash base64 rusak
	}

	for _, invalidHash := range invalidHashes {
		match := hasher.Compare(invalidHash, pass)
		assert.False(t, match, "Hash dengan format invalid '%s' seharusnya tidak cocok", invalidHash)
	}
}

// TestHash_CreatesDifferentHashes menguji bahwa salt bekerja
func TestHash_CreatesDifferentHashes(t *testing.T) {
	hasher := password.NewArgon2idHasher()
	pass := "password-yang-sama"

	// Hasilkan dua hash untuk password yang sama
	hash1, err1 := hasher.Hash(pass)
	hash2, err2 := hasher.Hash(pass)

	require.NoError(t, err1)
	require.NoError(t, err2)

	// Karena salt-nya acak, hash-nya harus berbeda
	assert.NotEqual(t, hash1, hash2, "Dua hash untuk password yang sama harus berbeda karena salt")

	// Tapi keduanya harus tetap valid
	assert.True(t, hasher.Compare(hash1, pass), "Hash 1 harus tetap valid")
	assert.True(t, hasher.Compare(hash2, pass), "Hash 2 harus tetap valid")
}
