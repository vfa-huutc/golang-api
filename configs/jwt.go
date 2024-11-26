package configs

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
)

type CustomClaims struct {
	Username             string `json:"username"` // Custom field
	jwt.RegisteredClaims        // // Embed standard claims
}

var jwtKey = []byte(utils.GetEnv("JWT_KEY", "replace_your_key"))

type JwtResult struct {
	Token     string
	ExpiresAt int64
}

func GenerateToken(username string) (*JwtResult, error) {
	expiresAt := jwt.NewNumericDate(time.Now().Add(time.Hour))
	claims := CustomClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expiresAt,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString(jwtKey)

	if err != nil {
		return nil, err
	}

	return &JwtResult{
		Token:     token,
		ExpiresAt: expiresAt.Unix(),
	}, nil
}

func ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err

}
