package database

import (
	"database/sql"
	"fmt"

	"github.com/AsishMandoi/iitk-coin/global"
	_ "github.com/mattn/go-sqlite3"
)

// Open the DB `iikusers.db` using sqlite, set _journal_mode to WAL
var db, cantOpenErr = sql.Open("sqlite3", "file:iitkusers.db?cache=shared&mode=rwc&_journal_mode=WAL")

// Open the DB `iikusers.db` using sqlite, _journal_mode is set to DELETE by default
// var db, cantOpenErr = sql.Open("sqlite3", "iitkusers.db")

// Ideally this function should be called before calling any other function in the `database` package
func Initialize() (string, error) {
	if cantOpenErr != nil {
		return "Could not access database", cantOpenErr
	}

	db.SetMaxOpenConns(1)

	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users(rollno INTEGER PRIMARY KEY not null, name TEXT not null, password TEXT not null, batch TEXT not null, coins DOUBLE not null);")
	if err != nil {
		return "Could not create table `users`", err
	}
	return "", nil
}

// Check if user with the given rollno is present in the DB
func FindRoll(rollno int) bool {
	row := db.QueryRow("SELECT COUNT(*) FROM users WHERE rollno=?;", rollno)
	var roll_present bool
	row.Scan(&roll_present)
	return roll_present
}

// Get the password for the given rollno (if present in the DB)
func GetPwd(rollno int) (string, error) {
	row := db.QueryRow("SELECT password FROM users WHERE rollno=?;", rollno)
	var pwd string
	err := row.Scan(&pwd)
	return pwd, err
}

// Get the batch for the given rollno (if present in the DB)
func GetBatch(rollno int) (string, error) {
	row := db.QueryRow("SELECT batch FROM users WHERE rollno=?;", rollno)
	var batch string
	err := row.Scan(&batch)
	return batch, err
}

// Add a user into the DB
func Add(usr global.Stu) (string, error) {
	// Check if user with the given rollno is present in the DB
	pres := FindRoll(usr.Rollno)

	// If such a user is already present, don't change anything in the DB.
	if pres {
		return "", fmt.Errorf("User already present")
	} else if _, err := db.Exec("INSERT INTO users(rollno, name, password, batch, coins) VALUES (?, ?, ?, ?, ?);", usr.Rollno, usr.Name, usr.Password, usr.Batch, 0); err != nil {
		return "Could not add user into the database", err
	}
	return "Added user successfully.", nil
}
