package server

import (
	"time"

	"github.com/AsishMandoi/iitk-coin/global"

	jwt "github.com/dgrijalva/jwt-go"
)

func GenJWT(usrRollno int) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	// Currently only the rollno and the expiration time are used as claims
	claims["rollno"] = usrRollno
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	// A token (tokenString) is generated using the header, the payload (using the claims) and the encoded secret key (gobal.SignatureKey)
	tokenString, err := token.SignedString(global.SignatureKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
