package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/sunnywalden/sync-data/pkg/models"
)

// GenerateToken, generate a token for platform user
func GenerateToken(user *models.PlatUser, authKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.UserName,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(authKey))
}

