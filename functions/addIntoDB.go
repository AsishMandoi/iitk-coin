package functions

import (
	"database/sql"
	"strconv"

	"github.com/AsishMandoi/iitk-coin/global"
)

func AddIntoDB(usr global.Stu) (string, error) {
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
