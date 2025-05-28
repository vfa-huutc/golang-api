package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
)

func TestHashPasswordAndCheckPasswordHash(t *testing.T) {
	password := "MyS3cr3tP@ssw0rd"

	// Hash the password
	hashed := utils.HashPassword(password)
	assert.NotEmpty(t, hashed, "Hashed password should not be empty")

	// Check that the correct password matches
	match := utils.CheckPasswordHash(password, hashed)
	assert.True(t, match, "The password should match the hash")

	// Check that a wrong password does not match
	wrongMatch := utils.CheckPasswordHash("wrongpassword", hashed)
	assert.False(t, wrongMatch, "A wrong password should not match the hash")

	// Check empty hash returns false
	assert.False(t, utils.CheckPasswordHash(password, ""), "Empty hash should not match any password")
}
