package database

import "database/sql"

// Open the DB `iikusers.db` using sqlite, set _journal_mode to WAL
var db, errNoDB = sql.Open("sqlite3", "file:iitkusers.db?cache=shared&mode=rwc&_journal_mode=WAL")

// Open the DB `iikusers.db` using sqlite, _journal_mode is set to DELETE by default
// var db, errNoDB = sql.Open("sqlite3", "iitkusers.db")

var InitMsg string
var InitErr error

func init() {
	if errNoDB != nil {
		InitMsg, InitErr = "Could not access database", errNoDB
		return
	}

	db.SetMaxOpenConns(1)

	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users(rollno INTEGER PRIMARY KEY not null, name TEXT not null, password TEXT not null, batch TEXT not null, role TEXT not null, coins DOUBLE not null);")
	if err != nil {
		InitMsg, InitErr = "Could not create table `users`", err
		return
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS transactions(id INTEGER PRIMARY KEY not null, type TEXT not null, sender INTEGER not null, receiver INTEGER not null, amount INTEGER not null, description TEXT not null, doneon TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP);")
	if err != nil {
		InitMsg, InitErr = "Could not create table `transactions`", err
	}
}
