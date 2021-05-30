package main

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type stu struct {
	rollno int
	name   string
}

func AddNewStudent(db *sql.DB, student stu) {
	// Checking if user with the given student.rollno is present in the DB
	row := db.QueryRow("SELECT COUNT(*) FROM users WHERE rollno=" + strconv.Itoa(student.rollno))
	var present bool
	row.Scan(&present)

	if !present {
		_, err := db.Exec("INSERT INTO users(rollno, name) VALUES (?, ?)", student.rollno, student.name)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Added user %v: %v successfully.\n", student.rollno, student.name)
	} else { // If such a user is already present, just display that message
		fmt.Printf("User with roll number %v is already present!\n", student.rollno)
	}
}

func main() {
	db, err := sql.Open("sqlite3", "./iitkusers.db") // A .db file is created in the same directory which is the database
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users(rollno INTEGER not null, name TEXT not null);")
	if err != nil {
		panic(err)
	}
	AddNewStudent(db, stu{190197, "Someone Cool"})
}
