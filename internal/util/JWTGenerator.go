package util

import (
	"AuthAPI/internal/services/auth"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateTokenFromGoogleIdentity(userInfo auth.GoogleUserInfo, jwtSecret string) (string, error) {
	claims := jwt.MapClaims{
		"userId": GenerateUUIDFromString(userInfo.ID),
		"email":  userInfo.Email,
		"Name":   userInfo.Name,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
