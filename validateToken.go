package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

func validateToken(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	payload := &secretpageRespBodyFormat{"-", "-", "-"} // Body of the response to be sent

	// The following function will be called when the signup function ends.
	defer func() {
		// Encode the payload (struct) into a json object and then send the json encoded body in the response.
		json.NewEncoder(w).Encode(*payload)
	}()
	// If the authorization token was not provided in the header--
	if r.Method == "GET" {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" || tokenString == "Bearer " {
			w.WriteHeader(http.StatusBadRequest)
			payload.Error = "Authorization token not found"
		} else {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				if t.Method.Alg() != "HS256" {
					return []byte(""), errors.New("invalid signing method")
				}
				return signatureKey, nil
			})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				payload.Error = err.Error()
				return
			} else if !token.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				payload.Error = "Invalid authorization token"
				return
			}
			payload.Message = "SUCESS"
			payload.Data = "Dummy data"
		}
	} else {
		payload.Message = "Welcome to /secretpage! Please use a GET request to get authorized."
	}
}
