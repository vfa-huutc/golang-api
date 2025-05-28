package utils

import (
	"testing"

	"github.com/vfa-khuongdv/golang-cms/internal/utils"
)

// TestGenerateRandomStringLength checks that the generated string is of the correct length
func TestGenerateRandomStringLength(t *testing.T) {
	length := 10
	randomStr := utils.GenerateRandomString(length)

	if len(randomStr) != length {
		t.Errorf("Expected length %d, but got %d", length, len(randomStr))
	}
}

// TestGenerateRandomStringUniqueness checks that consecutive calls return different strings (not guaranteed, but very likely)
func TestGenerateRandomStringUniqueness(t *testing.T) {
	str1 := utils.GenerateRandomString(10)
	str2 := utils.GenerateRandomString(10)

	if str1 == str2 {
		t.Errorf("Expected different strings but got the same: %s", str1)
	}
}

// TestGenerateRandomStringCharset ensures the result only contains allowed characters
func TestGenerateRandomStringCharset(t *testing.T) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	allowed := map[rune]bool{}
	for _, ch := range charset {
		allowed[ch] = true
	}

	randomStr := utils.GenerateRandomString(100)
	for _, ch := range randomStr {
		if !allowed[ch] {
			t.Errorf("Generated string contains invalid character: %c", ch)
		}
	}
}
