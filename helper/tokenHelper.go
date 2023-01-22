package helper

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

type TokenDetails struct {
	name string
	jwt.StandardClaims
}

func GenerateToken(userName string) (token string) {
	claims := &TokenDetails{
		name: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(time.Now().Day())).Unix(),
		},
	}

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	secretString := os.Getenv("SECRET_KEY")
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretString))
	if err != nil {
		log.Fatal(err)
	}

	return token
}
