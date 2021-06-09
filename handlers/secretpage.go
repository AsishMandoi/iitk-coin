package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/AsishMandoi/iitk-coin/global"
	jwt "github.com/dgrijalva/jwt-go"
)

// GET request format (in the header) -> --header
func ValidateToken(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	payload := &global.SecretpageRespBodyFormat{"-", "-", "-"} // Body of the response to be sent

	// The following function will be called when the signup function ends.
	defer func() {
		// Encode the payload (struct) into a json object and then send the json encoded body in the response.
		json.NewEncoder(w).Encode(*payload)
	}()

	// If the authorization token was not provided in the header--
	if r.Method == "GET" {

		tokenString := r.Header.Get("Authorization")

		// If the header `Authorization` is empty...
		if tokenString = strings.TrimPrefix(tokenString, "Bearer "); tokenString == "" {
			w.WriteHeader(http.StatusBadRequest)
			payload.Error = "Authorization token not found"
			return
		}

		// Get the (unverified) token from the header.
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != "HS256" {
				return []byte(""), errors.New("invalid signing method")
			}
			return global.SignatureKey, nil
		})

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			payload.Error = err.Error()
			return
		} else if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			payload.Error = "Invalid authorization token"
			return
		}

		// Since there is no error proper information from the secretpage is returned with the response.
		payload.Message = "SUCCESS"
		payload.Data = "Dummy data"
	} else {
		w.WriteHeader(http.StatusNotImplemented)
		payload.Message = "Welcome to /secretpage! Please use a GET request to get authorized."
	}
}
