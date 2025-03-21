package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey = contextKey("userID")
const AuthHeaderKey = contextKey("Authorization")

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateAccessToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString(jwtSecret)
}
