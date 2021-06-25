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

func Transact(
	Tx struct {
		Sender   int
		Receiver int
		Amount   float64
	},
	recAmt float64) error {

	txn, err := db.Begin()
	if err != nil {
		return err
	}
	defer txn.Rollback()

	res, err := txn.Exec("UPDATE users SET coins=coins-($1) WHERE rollno=($2) AND coins>=($1);", Tx.Amount, Tx.Sender)
	if err != nil {
		return fmt.Errorf("Could not send amount; %v", err)
	}
	if cntRows, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("Could not send amount; %v", err)
	} else if cntRows != 1 {
		return fmt.Errorf("Could not send amount; The sender may not have sufficient coins")
	}

	res, err = txn.Exec("UPDATE users SET coins=CASE WHEN coins+($1)<($2) THEN coins+($1) ELSE ($2) END WHERE rollno=$3;", recAmt, cap, Tx.Receiver)
	if err != nil {
		return fmt.Errorf("Could not recieve amount; %v", err)
	}
	if cntRows, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("Could not recieve amount; %v", err)
	} else if cntRows != 1 {
		return fmt.Errorf("Could not recieve amount; Possible error - Invalid receiver")
	}

	err = txn.Commit()
	if err != nil {
		return err
	}

	return nil
}

func Reward(Rwd struct {
	Receiver int     `json:"receiver"`
	Amount   float64 `json:"amount"`
}) error {
	txn, err := db.Begin()
	if err != nil {
		return err
	}
	defer txn.Rollback()

	res, err := txn.Exec("UPDATE users SET coins=CASE WHEN coins+($1)<($2) THEN coins+($1) ELSE ($2) END WHERE rollno=$3;", Rwd.Amount, cap, Rwd.Receiver)
	if err != nil {
		return fmt.Errorf("Could not recieve amount; %v", err)
	}
	if cntRows, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("Could not recieve amount; %v", err)
	} else if cntRows != 1 {
		return fmt.Errorf("Could not recieve amount; Possible error - Invalid receiver")
	}

	err = txn.Commit()
	if err != nil { // error while committing
		return err
	}

	return nil
}
