package utils

import (
	"crypto/rand"
	"io"
)

// GenerateOTP menghasilkan kode numerik acak dengan panjang yang ditentukan.
func GenerateOTP(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, buffer)
	if err != nil {
		return "", err
	}

	otpChars := "1234567890"
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%len(otpChars)]
	}

	return string(buffer), nil
}
