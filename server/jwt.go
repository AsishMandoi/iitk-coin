package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/AsishMandoi/iitk-coin/global"
	jwt "github.com/dgrijalva/jwt-go"
)

// Generate a JWT for a given user (rollno, batch, role)
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

// This authorizes a request (with a token) and returns the claims if successful else returns errors and a status code
func ValidateJWT(r *http.Request) (int, jwt.MapClaims, error) {

	tokenString := r.Header.Get("Authorization")

	// If the header `Authorization` is empty...
	if tokenString = strings.TrimPrefix(tokenString, "Bearer "); tokenString == "" {
		return 400, nil, fmt.Errorf("Authorization token not found")
	}

	// Get the (unverified) token from the header.
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != "HS256" {
			return []byte(""), fmt.Errorf("Invalid signing method")
		}
		return global.SignatureKey, nil
	})

	if err != nil {
		return 400, nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 401, nil, fmt.Errorf("Invalid authorization token")
	}
	return 200, claims, nil
}
