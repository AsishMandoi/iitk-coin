package server

import (
	"time"

	"github.com/AsishMandoi/iitk-coin/global"

	jwt "github.com/dgrijalva/jwt-go"
)

func GenJWT(usr struct {
	Rollno int
	Batch  string
	Role   string
}) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	// The rollno, batch, role ("Admin"/"CoreTeam"/"Coordinator"/nil) and the expiration time are used as claims
	claims["rollno"] = usr.Rollno
	claims["batch"] = usr.Batch
	claims["role"] = usr.Role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	// A token (tokenString) is generated using the header, the payload (using the claims) and the encoded secret key (gobal.SignatureKey)
	tokenString, err := token.SignedString(global.SignatureKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
