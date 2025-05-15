package utils

import (
	"math/rand"
	"time"
)

// GenerateRandomString generates a random string of specified length using alphanumeric characters
// Parameters:
//   - n: length of the random string to generate
//
// Returns:
//   - string: randomly generated alphanumeric string of length n
func GenerateRandomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	result := make([]byte, n)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result)
}
