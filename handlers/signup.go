package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AsishMandoi/iitk-coin/database"
	"github.com/AsishMandoi/iitk-coin/global"
	"github.com/AsishMandoi/iitk-coin/server"
	"golang.org/x/crypto/bcrypt"
)

func Signup(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var payload = &global.DefaultRespBody{} // Body of the response to be sent

	if r.Method == "POST" {
		// Converting the body into a json object
		var usr global.Stu
		err := json.NewDecoder(r.Body).Decode(&usr)
		if err != nil {
			server.Respond(w, payload, 400, "Could not decode body of the request", err.Error())
			return
		}

		// Hashing the password with a cost of 10
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usr.Password), 10)
		if err != nil {
			server.Respond(w, payload, 500, "Could not generate hash from password", err.Error())
			return
		}

		// The hashed user password should be stored instead of the plaintext version
		usr.Password = string(hashedPassword)

		// Handle initialization errors in SQLite DB
		if msg, err := database.InitMsg, database.InitErr; err != nil {
			server.Respond(w, payload, 500, msg, err.Error())
			return
		}

		// Add user into DB
		msg, err := database.Add(usr)
		if err != nil {
			if err.Error() == "User already present" {
				server.Respond(w, payload, 409, "You have alrady signed up!", err.Error())
			} else {
				server.Respond(w, payload, 500, msg, err.Error())
			}
		} else {
			server.Respond(w, payload, 201, msg, nil)
		}
	} else {
		server.Respond(w, payload, 501, "Welcome to signup page! Please use a POST method to signup.", nil)
	}
}
