package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword takes a plain text password string and returns a hashed version using bcrypt
// It uses the default cost factor for the hashing algorithm
// Returns the hashed password as a string and any error that occurred during hashing
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}
