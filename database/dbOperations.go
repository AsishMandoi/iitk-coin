package database

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/AsishMandoi/iitk-coin/global"
	_ "github.com/mattn/go-sqlite3"
)

// Open the DB `iikusers.db` using sqlite
var db, cantOpenErr = sql.Open("sqlite3", "./iitkusers.db")

// Ideally this function should be called before calling any other function in the `database` package
func Initialize() (string, error) {
	if cantOpenErr != nil {
		return "Could not access database", cantOpenErr
	}
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users(rollno INTEGER not null, name TEXT not null, password TEXT not null);")
	if err != nil {
		return "Could not create table `users`", err
	}
	return "", nil
}

// Check if user with the given rollno is present in the DB
func CheckTableForRoll(rollno int) bool {
	row := db.QueryRow("SELECT COUNT(*) FROM users WHERE rollno=" + strconv.Itoa(rollno))
	var roll_present bool
	row.Scan(&roll_present)
	return roll_present
}

// Get the password for the given rollno (if present in the DB)
func GetPwdForRoll(rollno int) (string, error) {
	row := db.QueryRow("SELECT password FROM users WHERE rollno=" + strconv.Itoa(rollno))
	var pwd string
	err := row.Scan(&pwd)
	return pwd, err
}

// Add a user into the DB
func Add(usr global.Stu) (string, error) {
	// Check if user with the given rollno is present in the DB
	pres := CheckTableForRoll(usr.Rollno)

	// If such a user is already present, don't change anything in the DB.
	if pres {
		return "", errors.New("User already present")
	} else if _, err := db.Exec("INSERT INTO users(rollno, name, password) VALUES (?, ?, ?)", usr.Rollno, usr.Name, usr.Password); err != nil {
		return "Could not add user into the database", err
	}
	return "Added user successfully.", nil
}
