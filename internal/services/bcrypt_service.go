package services

import "golang.org/x/crypto/bcrypt"

type IBcryptService interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hashPassword string) bool
	HashPasswordWithCost(password string, cost int) (string, error)
}

type BcryptService struct{}

func NewBcryptService() IBcryptService {
	return &BcryptService{}
}

// HashPassword hashes a password using bcrypt with the default cost
// Returns the hashed password as a string, or an error if hashing fails
func (s *BcryptService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash compares a plain text password with a hashed password
// Returns true if they match, false otherwise
func (s *BcryptService) CheckPasswordHash(password, hashPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}

// HashPasswordWithCost hashes a password using bcrypt with a specified cost
// Returns the hashed password as a string, or an error if hashing fails
func (s *BcryptService) HashPasswordWithCost(password string, cost int) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
