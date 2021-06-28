package database

import (
	"fmt"

	"github.com/AsishMandoi/iitk-coin/global"
)

var cap = global.MaxCap

// Get the balance for the given rollno (if present in the DB)
func GetCoins(rollno interface{}) (float64, error) {
	row := db.QueryRow("SELECT coins FROM users WHERE rollno=?;", rollno)
	var coins float64
	err := row.Scan(&coins)
	return coins, err
}

// Lets an authorized, eligible user send coins to another valid, eligible user
func Transact(Tx global.TxnBody, amtRcvd float64) (interface{}, error) {

	txn, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer txn.Rollback()

	res, err := txn.Exec("UPDATE users SET coins=coins-($1) WHERE rollno=($2) AND coins>=($1);", Tx.Amount, Tx.Sender)
	if err != nil {
		return nil, fmt.Errorf("Could not send amount; %v", err)
	} else if cntRows, err := res.RowsAffected(); err != nil {
		return nil, fmt.Errorf("Could not send amount; %v", err)
	} else if cntRows != 1 {
		return nil, fmt.Errorf("Could not send amount; The sender may not have sufficient coins")
	}

	res, err = txn.Exec("UPDATE users SET coins=CASE WHEN coins+($1)<($2) THEN coins+($1) ELSE ($2) END WHERE rollno=($3);", amtRcvd, cap, Tx.Receiver)
	if err != nil {
		return nil, fmt.Errorf("Could not recieve amount; %v", err)
	} else if cntRows, err := res.RowsAffected(); err != nil {
		return nil, fmt.Errorf("Could not recieve amount; %v", err)
	} else if cntRows != 1 {
		return nil, fmt.Errorf("Could not recieve amount; Possible error - Invalid receiver")
	}

	// Storing the transaction as a log in another table
	res, err = txn.Exec("INSERT INTO transactions(type, sender, receiver, amount, description) VALUES (?, ?, ?, ?, ?)", "Transfer", Tx.Sender, Tx.Receiver, Tx.Amount, Tx.Descr)
	if err != nil {
		return nil, fmt.Errorf("Could not log transfer; %v", err)
	} else if cntRows, err := res.RowsAffected(); err != nil {
		return nil, fmt.Errorf("Could not log transfer; %v", err)
	} else if cntRows != 1 {
		return nil, fmt.Errorf("Could not log transfer")
	}

	txid, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("Could not log transfer; %v", err)
	}

	err = txn.Commit()
	if err != nil {
		return nil, err
	}

	return txid, nil
}

// Lets an authorized Admin reward coins to a valid user
func Reward(Tx global.TxnBody) (interface{}, error) {
	txn, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer txn.Rollback()

	res, err := txn.Exec("UPDATE users SET coins=CASE WHEN coins+($1)<($2) THEN coins+($1) ELSE ($2) END WHERE rollno=($3);", Tx.Amount, cap, Tx.Receiver)
	if err != nil {
		return nil, fmt.Errorf("Could not recieve amount; %v", err)
	} else if cntRows, err := res.RowsAffected(); err != nil {
		return nil, fmt.Errorf("Could not recieve amount; %v", err)
	} else if cntRows != 1 {
		return nil, fmt.Errorf("Could not recieve amount; Possible error - Invalid receiver")
	}

	// Storing the transaction as a log in another table
	res, err = txn.Exec("INSERT INTO transactions(type, sender, receiver, amount, description) VALUES (?, ?, ?, ?, ?)", "Reward", Tx.Sender, Tx.Receiver, Tx.Amount, Tx.Descr)
	if err != nil {
		return nil, fmt.Errorf("Could not log reward; %v", err)
	} else if cntRows, err := res.RowsAffected(); err != nil {
		return nil, fmt.Errorf("Could not log reward; %v", err)
	} else if cntRows != 1 {
		return nil, fmt.Errorf("Could not log reward")
	}

	txid, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("Could not log reward; %v", err)
	}

	err = txn.Commit()
	if err != nil {
		return nil, err
	}

	return txid, nil
}
