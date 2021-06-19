package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/AsishMandoi/iitk-coin/global"
	jwt "github.com/dgrijalva/jwt-go"
)

// This authorizes a request (with a token) and returns the claims if successful else returns errors and a status code
func Authorize(r *http.Request) (int, jwt.MapClaims, error) {

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
