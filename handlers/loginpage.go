package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AsishMandoi/iitk-coin/database"
	"github.com/AsishMandoi/iitk-coin/global"
	"github.com/AsishMandoi/iitk-coin/server"
	"golang.org/x/crypto/bcrypt"
)

// POST request format (in the body) -> {"rollno": 190197, "password": "sTr0nG-p@$5w0rD"}
func Login(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	payload := &global.LoginRespBodyFormat{} // Body of the response to be sent

	if r.Method == "POST" {
		// Converting the body into a json object
		var usr global.LoginInputFormat
		if err := json.NewDecoder(r.Body).Decode(&usr); err != nil {
			server.Respond(w, payload, 400, "Could not decode body of the request", err.Error(), "-")
			return
		}

		if msg, err := database.Initialize(); err != nil {
			server.Respond(w, payload, 500, msg, err.Error(), "-")
			return
		}

		pwd, err := database.GetPwdForRoll(usr.Rollno)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				server.Respond(w, payload, 400, "Could not identify user with given roll no; Please make sure that you have signed up or check your roll no again.", err.Error(), "-")
				return
			}
			server.Respond(w, payload, 500, "Something went wrong :(", err.Error(), "-")
			return
		}

		// Comparing the hash generated from the entered password with the bcrypt-hash (stored in the DB)
		if err := bcrypt.CompareHashAndPassword([]byte(pwd), []byte(usr.Password)); err != nil {
			server.Respond(w, payload, 401, "Login unsuccessful; Invalid user credentials", err.Error(), "-")
		} else {
			payload.Message = "Login successful"
			if token, err := server.GenJWT(usr); err != nil {
				server.Respond(w, payload, 403, "Login successful; Token could not be generated", err.Error(), "-")
			} else {
				server.Respond(w, payload, 200, "Login successful; Token generated successfully", "-", token)
			}
		}
	} else {
		server.Respond(w, payload, 501, "Welcome to login page! Please use a POST request to login.", "-", "-")
	}
}
