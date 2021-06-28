package database

import (
	"fmt"

	"github.com/AsishMandoi/iitk-coin/global"
	_ "github.com/mattn/go-sqlite3"
)

// Get the user details (password, batch, and role) for the given rollno (if present)
func GetUsrDetails(rollno int) (struct{ Pwd, Batch, Role string }, error) {
	row := db.QueryRow("SELECT password, batch, role FROM users WHERE rollno=?;", rollno)
	var usrDetails struct{ Pwd, Batch, Role string }
	err := row.Scan(&usrDetails.Pwd, &usrDetails.Batch, &usrDetails.Role)
	return usrDetails, err
}

// Get the batch and role for the given rollno (if present)
func GetBatchnRole(rollno int) (string, string, error) {
	row := db.QueryRow("SELECT batch, role FROM users WHERE rollno=?;", rollno)
	var batch, role string
	err := row.Scan(&batch, &role)
	return batch, role, err
}

// Get the count of the number of events in which the given rollno (if present) has participated in
func GetCntEvents(rollno int) (int, error) {
	row := db.QueryRow("SELECT COUNT(*) FROM transactions WHERE type='Reward' AND receiver=?", rollno)
	var count int
	err := row.Scan(&count)
	return count, err
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
