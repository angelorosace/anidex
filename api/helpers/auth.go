package helpers

import (
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
)

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid signing method")
		}
		return []byte(os.Getenv("SALT")), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
