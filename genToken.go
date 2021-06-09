package main

import (
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var signatureKey = []byte(os.Getenv("SECRET_KEY"))

func genJWT(payload loginInputFormat) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["rollno"] = payload.Rollno
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()

	tokenString, err := token.SignedString(signatureKey)
	if err != nil {
		// fmt.Errorf("Something went wrong: %s", err.Error())
		return "", err
	}

	// fmt.Println(tokenString)
	return tokenString, nil
}
