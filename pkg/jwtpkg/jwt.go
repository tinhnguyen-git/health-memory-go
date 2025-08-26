package jwtpkg

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokenHS256(secret string, subject string, email string, provider string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":   subject,
		"email": email,
		"prov":  provider,
		"exp":   time.Now().Add(ttl).Unix(),
		"iat":   time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}
