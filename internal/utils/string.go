package utils

import (
	"math/rand"
	"time"
)

// GenerateRandomString generates a random string of length n using the given character set
func GenerateRandomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return generateStringWithCharset(n, charset)
}

// generateStringWithCharset generates a random string of length n using a custom charset
func generateStringWithCharset(n int, charset string) string {
	// Create a new random generator with a random seed
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	result := make([]byte, n)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result)
}
