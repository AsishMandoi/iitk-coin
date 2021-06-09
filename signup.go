package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func addIntoDB(usr Stu) (string, error) {
	// Open the DB `iikusers.db`
	db, err := sql.Open("sqlite3", "./iitkusers.db")
	if err != nil {
		return "", err
	}
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS users(rollno INTEGER not null, name TEXT not null, password TEXT not null);"); err != nil {
		return "", err
	}
	// Check if user with the given rollno is present in the DB
	row := db.QueryRow("SELECT COUNT(*) FROM users WHERE rollno=" + strconv.Itoa(usr.Rollno))
	var present bool
	row.Scan(&present)

	// If such a user is already present, don't change anything in the DB.
	if present {
		return "User with roll number " + strconv.Itoa(usr.Rollno) + " is already present!", nil
	}

	if _, err := db.Exec("INSERT INTO users(rollno, name, password) VALUES (?, ?, ?)", usr.Rollno, usr.Name, usr.Password); err != nil {
		return "", err
	}
	return "Added user successfully.", nil
}

func signup(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	payload := &signupRespBodyFormat{"-", "-"}	// Body of the response to be sent

	// The following function will be called when the signup function ends.
	defer func() {
		// Encode the payload (struct) into a json object and then send the json encoded body in the response.
		json.NewEncoder(w).Encode(*payload)
	}()

	if r.Method == "POST" {
		// Converting the body into a json object
		var usr Stu
		err := json.NewDecoder(r.Body).Decode(&usr)
		if err != nil {
			payload.Message = "Something went wrong :("
			payload.Error = err.Error()
			return
		}

		// Hashing the password with a cost of 10
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usr.Password), 10)
		if err != nil {
			payload.Message = "Something went wrong :("
			payload.Error = err.Error()
			return
		}

		// The hashed user password should be stored instead of the plaintext version
		usr.Password = string(hashedPassword)

		msg, err := addIntoDB(usr)
		if err != nil {
			payload.Message = "Something went wrong :("
			payload.Error = err.Error()
		} else {
			payload.Message = msg
		}
	} else {
		payload.Message = "Welcome to signup page! Please use a POST request to signup."
	}
}
