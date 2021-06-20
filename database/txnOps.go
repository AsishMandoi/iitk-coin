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

func Transact(Tx struct {
	Sender   int     `json:"sender"`
	Receiver int     `json:"receiver"`
	Amount   float64 `json:"amount"`
}, recAmt float64) error {

	txn, err := db.Begin()
	if err != nil {
		return err
	}
	defer txn.Rollback()

	res, err := txn.Exec("UPDATE users SET coins=coins-(?) WHERE rollno=(?) AND coins>=(?);", Tx.Amount, Tx.Sender, Tx.Amount)
	if err != nil {
		return fmt.Errorf("Could not send amount; %v", err)
	}
	if cntRows, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("Could not send amount; %v", err)
	} else if cntRows != 1 {
		return fmt.Errorf("Could not send amount; Possible errors - i) Invalid sender, ii) Insufficient coins")
	}

	res, err = txn.Exec("UPDATE users SET coins=CASE WHEN coins+(?)<(?) THEN coins+(?) ELSE (?) END WHERE rollno=?;", recAmt, cap, recAmt, cap, Tx.Receiver)
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

	res, err := txn.Exec("UPDATE users SET coins=CASE WHEN coins+(?)<(?) THEN coins+(?) ELSE (?) END WHERE rollno=?;", Rwd.Amount, cap, Rwd.Amount, cap, Rwd.Receiver)
	if err != nil {
		return fmt.Errorf("Could not recieve amount; %v", err)
	}
	if cntRows, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("Could not recieve amount; %v", err)
	} else if cntRows != 1 {
		// Try to handle specific errors here if possible
		return fmt.Errorf("Could not recieve amount; Possible error - Invalid receiver")
	}

	err = txn.Commit()
	if err != nil { // error while committing
		return err
	}

	return nil
}
