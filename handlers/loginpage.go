package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/AsishMandoi/iitk-coin/functions"
	"github.com/AsishMandoi/iitk-coin/global"
	"golang.org/x/crypto/bcrypt"
)

// POST request format (in the body) -> {"rollno": 190197, "password": "sTr0nG-p@$5w0rD"}
func Login(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	payload := &global.LoginRespBodyFormat{"-", "-", "-"} // Body of the response to be sent

	// The following function will be called when the login function ends.
	defer func() {
		// Encode the payload (struct) into a json object and then send the json encoded body in the response.
		json.NewEncoder(w).Encode(*payload)
	}()

	if r.Method == "POST" {
		// Converting the body into a json object
		var usr global.LoginInputFormat
		if err := json.NewDecoder(r.Body).Decode(&usr); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			payload.Message = "Could not decode body of the request"
			payload.Error = err.Error()
			return
		}

		// Open the DB `iikusers.db`
		db, err := sql.Open("sqlite3", "./iitkusers.db")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			payload.Message = "Could not open database"
			payload.Error = err.Error()
			return
		}

		// Check if user with the given rollno is present in the DB
		row := db.QueryRow("SELECT password FROM users WHERE rollno=" + strconv.Itoa(usr.Rollno))
		var pwd []byte
		if err := row.Scan(&pwd); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			payload.Message = "Something went wrong :("
			payload.Error = err.Error()
			return
		}

		// Comparing the hash generated from the entered password with the bcrypt-hash (stored in the DB)
		if err := bcrypt.CompareHashAndPassword(pwd, []byte(usr.Password)); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			payload.Message = "Login unsuccessful. Invalid user credentials"
			payload.Error = err.Error()
		} else {
			payload.Message = "Login successful"
			if token, err := functions.GenJWT(usr); err != nil {
				w.WriteHeader(http.StatusForbidden)
				payload.Message = "Token could not be generated."
				payload.Error = err.Error()
			} else {
				payload.Message += "; Token generated successfully"
				payload.Token = token
			}
		}
	} else {
		w.WriteHeader(http.StatusNotImplemented)
		payload.Message = "Welcome to login page! Please use a POST request to login."
	}
}
