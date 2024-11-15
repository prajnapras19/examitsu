package lib

import (
	"crypto/rand"
	"io"
)

func GenerateRandomString(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	randomData := make([]byte, length)

	// Open /dev/urandom
	_, err := io.ReadFull(rand.Reader, randomData)
	if err != nil {
		return "", err
	}

	// Map bytes to charset
	for i, b := range randomData {
		result[i] = charset[b%byte(len(charset))]
	}

	return string(result), nil
}
