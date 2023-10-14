package lib

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomToken() (string, error) {
	length := 32

	// Calculate the byte size needed for the given length
	byteSize := length / 4 * 3
	if length%4 > 0 {
		byteSize += 3
	}

	// Generate a random byte slice
	randomBytes := make([]byte, byteSize)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Encode the random bytes to a base64 string
	token := base64.URLEncoding.EncodeToString(randomBytes)

	// Trim '=' padding characters
	token = token[:length]

	return token, nil
}
