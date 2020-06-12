package apis

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

func TokenGenerator(platUserId string) (token *jwt.Token, err error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": platUserId,
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		return nil, err
	}

	log.Infof(tokenString, err)
	return tokenString, nil
}
