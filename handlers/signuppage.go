package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AsishMandoi/iitk-coin/functions"
	"github.com/AsishMandoi/iitk-coin/global"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// POST request format (in the body) -> {"rollno": 190197, "name": "Someone Cool", "password": "sTr0nG-p@$5w0rD"}
func Signup(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	payload := &global.SignupRespBodyFormat{"-", "-"} // Body of the response to be sent

	// The following function will be called when the signup function ends.
	defer func() {
		// Encode the payload (struct) into a json object and then send the json encoded body in the response.
		json.NewEncoder(w).Encode(*payload)
	}()

	if r.Method == "POST" {
		// Converting the body into a json object
		var usr global.Stu
		err := json.NewDecoder(r.Body).Decode(&usr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			payload.Message = "Something went wrong :("
			payload.Error = err.Error()
			return
		}

		// Hashing the password with a cost of 10
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usr.Password), 10)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			payload.Message = "Could not generate hash from password"
			payload.Error = err.Error()
			return
		}

		// The hashed user password should be stored instead of the plaintext version
		usr.Password = string(hashedPassword)

		msg, err := functions.AddIntoDB(usr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			payload.Message = "Could not add user into the database"
			payload.Error = err.Error()
		} else {
			payload.Message = msg
		}
	} else {
		w.WriteHeader(http.StatusNotImplemented)
		payload.Message = "Welcome to signup page! Please use a POST request to signup."
	}
}
