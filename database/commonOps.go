package database

import (
	"database/sql"
	"fmt"

	"github.com/AsishMandoi/iitk-coin/global"
	_ "github.com/mattn/go-sqlite3"
)

// Open the DB `iikusers.db` using sqlite, set _journal_mode to WAL
var db, errNoDb = sql.Open("sqlite3", "file:iitkusers.db?cache=shared&mode=rwc&_journal_mode=WAL")

// Open the DB `iikusers.db` using sqlite, _journal_mode is set to DELETE by default
// var db, errNoDb = sql.Open("sqlite3", "iitkusers.db")

// Ideally this function should be called before calling any other function in the `database` package
func Initialize() (string, error) {
	if errNoDb != nil {
		return "Could not access database", errNoDb
	}

	db.SetMaxOpenConns(1)

	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users(rollno INTEGER PRIMARY KEY not null, name TEXT not null, password TEXT not null, batch TEXT not null, role TEXT not null, coins DOUBLE not null);")
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
func GetUsrDetails(rollno int) (struct{ Pwd, Batch, Role string }, error) {
	row := db.QueryRow("SELECT password, batch, role FROM users WHERE rollno=?;", rollno)
	var usrDetails struct{ Pwd, Batch, Role string }
	err := row.Scan(&usrDetails.Pwd, &usrDetails.Batch, &usrDetails.Role)
	return usrDetails, err
}

// Get the batch for the given rollno (if present in the DB)
func GetBatchnRole(rollno int) (string, string, error) {
	row := db.QueryRow("SELECT batch, role FROM users WHERE rollno=?;", rollno)
	var batch, role string
	err := row.Scan(&batch, &role)
	return batch, role, err
}

// Add a user into the DB
func Add(usr global.Stu) (string, error) {
	_, err := db.Exec("INSERT INTO users(rollno, name, password, batch, role, coins) VALUES (?, ?, ?, ?, ?, ?);", usr.Rollno, usr.Name, usr.Password, usr.Batch, "", 0)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.rollno" {
			return "Could not add user into the database", fmt.Errorf("User #%v already present", usr.Rollno)
		}
		return "Could not add user into the database", err
	}
	return "Added user successfully.", nil
}
