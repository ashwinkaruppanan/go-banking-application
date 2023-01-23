package helper

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenDetails struct {
	UserType string
	Uid      string
	Name     string
	jwt.StandardClaims
}

var secretString string = os.Getenv("SECRET_KEY")

func GenerateToken(userName string, UserID string, UserType string) (token string) {
	claims := &TokenDetails{
		UserType: UserType,
		Uid:      UserID,
		Name:     userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretString))
	if err != nil {
		log.Panic(err)
	}

	return token
}

func ValidateToken(token string) (claims *TokenDetails, msg string) {
	tkn, err := jwt.ParseWithClaims(token, &TokenDetails{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretString), nil
	})

	if err != nil {
		log.Panic(err)
	}

	claims, ok := tkn.Claims.(*TokenDetails)
	if !ok {
		msg = "invalid token"
		return
	}

	if claims.ExpiresAt < time.Now().Unix() {
		msg = "token expird"
		return
	}

	return claims, msg
}
