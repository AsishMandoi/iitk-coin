package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/AsishMandoi/iitk-coin/global"
	"github.com/AsishMandoi/iitk-coin/server"
	jwt "github.com/dgrijalva/jwt-go"
)

// GET request format (in the header) -> --header "Authorization: Bearer <access token>"
func ValidateToken(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	payload := &global.SecretpageRespBodyFormat{} // Body of the response to be sent

	if r.Method == "GET" {

		tokenString := r.Header.Get("Authorization")

		// If the header `Authorization` is empty...
		if tokenString = strings.TrimPrefix(tokenString, "Bearer "); tokenString == "" {
			server.Respond(w, payload, 400, "-", "Authorization token not found", "-")
			return
		}

		// Get the (unverified) token from the header.
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != "HS256" {
				return []byte(""), errors.New("Invalid signing method")
			}
			return global.SignatureKey, nil
		})

		if err != nil {
			server.Respond(w, payload, 400, "-", err.Error(), "-")
			return
		} else if !token.Valid {
			server.Respond(w, payload, 401, "-", "Invalid authorization token", "-")
			return
		}

		// Since there are no more errors, the secretpage is responds with the confidential information.
		server.Respond(w, payload, 200, "SUCCESS", "-", "Dummy data")
	} else {
		server.Respond(w, payload, 501, "Welcome to /secretpage! Please use a GET request to get authorized.", "-", "-")
	}
}
