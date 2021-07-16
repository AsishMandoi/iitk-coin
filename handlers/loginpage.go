package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AsishMandoi/iitk-coin/database"
	"github.com/AsishMandoi/iitk-coin/global"
	"github.com/AsishMandoi/iitk-coin/server"
	"golang.org/x/crypto/bcrypt"
)

// POST request format (in the body) -> {"rollno": 190197, "password": "sTr0nG-p@$5w0rD"}
func Login(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	payload := &global.LoginRespBody{} // Body of the response to be sent

	if r.Method == "POST" {
		usr := struct {
			Rollno int    `json:"rollno"`
			Pwd    string `json:"password"`
		}{}

		// Converting the body into a json object
		if err := json.NewDecoder(r.Body).Decode(&usr); err != nil {
			server.Respond(w, payload, 400, "Could not decode body of the request", err.Error(), nil)
			return
		}

		// Handle initialization errors in DB
		if msg, err := database.InitMsg, database.InitErr; err != nil {
			server.Respond(w, payload, 500, msg, err.Error(), nil)
			return
		}

		usrDB, err := database.GetUsrDetails(usr.Rollno)
		if err != nil {
			if err == sql.ErrNoRows {
				server.Respond(w, payload, 400, fmt.Sprintf("Could not identify user with roll no %v", usr.Rollno), err.Error(), nil)
			} else {
				server.Respond(w, payload, 500, "Could not fetch user details", err.Error(), nil)
			}
			return
		}

		// Comparing the hash generated from the entered password with the bcrypt-hash (stored in the DB)
		if err := bcrypt.CompareHashAndPassword([]byte(usrDB.Pwd), []byte(usr.Pwd)); err != nil {
			server.Respond(w, payload, 401, "Login unsuccessful; Invalid user credentials", err.Error(), nil)
			return
		}

		if token, err := server.GenJWT(struct {
			Rollno      int
			Batch, Role string
		}{usr.Rollno, usrDB.Batch, usrDB.Role}); err != nil {
			server.Respond(w, payload, 403, "Login successful; Token could not be generated", err.Error(), nil)
		} else {
			server.Respond(w, payload, 200, "Login successful; Token generated successfully", nil, token)
		}
	} else {
		server.Respond(w, payload, 501, "Welcome to login page! Please use a POST method to login.", nil, nil)
	}
}
