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

		// Handle initialization errors in SQLite DB
		if msg, err := database.InitMsg, database.InitErr; err != nil {
			server.Respond(w, payload, 500, msg, err.Error(), nil)
			return
		}

		usrDB, err := database.GetUsrDetails(usr.Rollno)
		if err == sql.ErrNoRows {
			server.Respond(w, payload, 400, fmt.Sprintf("Could not identify user with roll no %v", usr.Rollno), err.Error(), nil)
			return
		} else if err != nil {
			server.Respond(w, payload, 500, "Could not fetch user details", err.Error(), nil)
			return
		}

		// Comparing the hash generated from the entered password with the bcrypt-hash (stored in the DB)
		if err := bcrypt.CompareHashAndPassword([]byte(usrDB.Pwd), []byte(usr.Pwd)); err == bcrypt.ErrMismatchedHashAndPassword {
			server.Respond(w, payload, 401, "Login unsuccessful: invalid user credentials", err.Error(), nil)
			return
		} else if err != nil {
			server.Respond(w, payload, 500, "Login unsuccessful: Failed to comapare hash and password", err.Error(), nil)
			return
		}

		if token, err := server.GenJWT(struct {
			Rollno             int
			Email, Batch, Role string
		}{usr.Rollno, usrDB.Email, usrDB.Batch, usrDB.Role}); err != nil {
			server.Respond(w, payload, 502, "Login successful; Token could not be generated", err.Error(), nil)
		} else {
			server.Respond(w, payload, 200, "Login successful; Token generated successfully", nil, token)
		}
	} else {
		server.Respond(w, payload, 501, "Welcome to login page! Please use a POST method to login.", nil, nil)
	}
}

func ResetPwd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	payload := &global.DefaultRespBody{} // Body of the response to be sent

	if r.Method == "POST" {

		body := struct {
			Otp    bool   `json:"send_otp"`
			OldPwd string `json:"old_password"`
			NewPwd string `json:"new_password"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			server.Respond(w, payload, 400, "Could not decode body of the request", err.Error())
			return
		}

		// Authorizing the request and obtaining the user's roll no
		statusCode, claims, err := server.ValidateJWT(r)
		if err != nil {
			server.Respond(w, payload, statusCode, nil, err.Error())
			return
		}

		rollno := int(claims["rollno"].(float64))

		if !body.Otp {
			usrDB, err := database.GetUsrDetails(rollno)
			if err == sql.ErrNoRows {
				server.Respond(w, payload, 400, fmt.Sprintf("Could not identify user with roll no %v", rollno), err.Error(), nil)
				return
			} else if err != nil {
				server.Respond(w, payload, 500, "Could not fetch user details", err.Error(), nil)
				return
			}

			// Comparing the hash generated from the entered password with the bcrypt-hash (stored in the DB)
			if err := bcrypt.CompareHashAndPassword([]byte(usrDB.Pwd), []byte(body.OldPwd)); err == bcrypt.ErrMismatchedHashAndPassword {
				server.Respond(w, payload, 401, "Invalid user credentials", err.Error(), nil)
				return
			} else if err != nil {
				server.Respond(w, payload, 500, "Failed to comapare hash and password", err.Error(), nil)
				return
			}
		}

		// Hashing the new password with a cost of 10
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.NewPwd), 10)
		if err != nil {
			server.Respond(w, payload, 500, "Could not generate hash from password", err.Error())
			return
		}

		// The hashed user password should be stored instead of the plaintext version
		body.NewPwd = string(hashedPassword)

		if body.Otp {
			email := claims["email"].(string)

			// Generate OTP, save it (along with other details) in the redis database with an expiry time and then send it
			if msg, err := server.SendOTP(email, global.PwdResetObj{rollno, body.NewPwd, ""}, "resetPwd"); err != nil {
				server.Respond(w, payload, 500, msg, err.Error())
				return
			}
			server.Respond(w, payload, 200, "Post your otp on http://localhost:8080/reset_password/confirm to confirm your transaction", nil)
			return
		}

		// Handle initialization errors in SQLite DB
		if msg, err := database.InitMsg, database.InitErr; err != nil {
			server.Respond(w, payload, 500, msg, err.Error())
			return
		}
		database.UpdPwd(rollno, body.NewPwd)
		server.Respond(w, payload, 200, "Password reset successful", nil)
	} else {
		server.Respond(w, payload, 501, "Welcome to /reset_password page! Please use a POST method to send coins to another user.", nil)
	}
}
