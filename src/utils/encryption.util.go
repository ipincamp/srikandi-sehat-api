package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

// EncryptToken encrypts a token with a password using AES-256-GCM
func EncryptToken(token, password string) (string, error) {
	// Derive a key from the password using SHA-256
	key := sha256.Sum256([]byte(password))

	// Create AES cipher
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the token
	ciphertext := gcm.Seal(nonce, nonce, []byte(token), nil)

	// Encode to base64 for safe transmission
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// DecryptToken decrypts an encrypted token with a password
func DecryptToken(encryptedToken, password string) (string, error) {
	// Derive the same key from the password
	key := sha256.Sum256([]byte(password))

	// Decode from base64
	ciphertext, err := base64.URLEncoding.DecodeString(encryptedToken)
	if err != nil {
		return "", err
	}

	// Create AES cipher
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Extract nonce
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
