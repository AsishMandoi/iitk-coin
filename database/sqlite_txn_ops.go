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
func Transact(Tx global.TxnObj) (interface{}, error) {

	txn, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer txn.Rollback()

	res, err := txn.Exec("UPDATE users SET coins=coins-($1) WHERE rollno=($2) AND coins>=($1);", Tx.AmtSent, Tx.Sender)
	if err != nil {
		return nil, fmt.Errorf("Could not send amount; %v", err)
	} else if cntRows, err := res.RowsAffected(); err != nil {
		return nil, fmt.Errorf("Could not send amount; %v", err)
	} else if cntRows != 1 {
		return nil, fmt.Errorf("Could not send amount; The sender may not have sufficient coins")
	}

	res, err = txn.Exec("UPDATE users SET coins=CASE WHEN coins+($1)<($2) THEN coins+($1) ELSE ($2) END WHERE rollno=($3);", Tx.AmtRcvd, cap, Tx.Receiver)
	if err != nil {
		return nil, fmt.Errorf("Could not receive amount; %v", err)
	} else if cntRows, err := res.RowsAffected(); err != nil {
		return nil, fmt.Errorf("Could not receive amount; %v", err)
	} else if cntRows != 1 {
		return nil, fmt.Errorf("Could not receive amount; Possible error - Invalid receiver")
	}

	// Storing the transaction as a log in another table
	res, err = txn.Exec("INSERT INTO transactions(type, sender, receiver, amount, description) VALUES (?, ?, ?, ?, ?)", "Transfer", Tx.Sender, Tx.Receiver, Tx.AmtSent, Tx.Descr)
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
func Reward(Tx global.TxnObj) (interface{}, error) {
	txn, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer txn.Rollback()

	res, err := txn.Exec("UPDATE users SET coins=CASE WHEN coins+($1)<($2) THEN coins+($1) ELSE ($2) END WHERE rollno=($3);", Tx.AmtSent, cap, Tx.Receiver)
	if err != nil {
		return nil, fmt.Errorf("Could not receive amount; %v", err)
	} else if cntRows, err := res.RowsAffected(); err != nil {
		return nil, fmt.Errorf("Could not receive amount; %v", err)
	} else if cntRows != 1 {
		return nil, fmt.Errorf("Could not receive amount; Possible error - Invalid receiver")
	}

	// Storing the transaction as a log in another table
	res, err = txn.Exec("INSERT INTO transactions(type, sender, receiver, amount, description) VALUES (?, ?, ?, ?, ?)", "Reward", Tx.Sender, Tx.Receiver, Tx.AmtSent, Tx.Descr)
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

// Update redeem status
func UpdRdmSts(redeemReq global.RedeemStatusUPDBody) (interface{}, error) {
	if redeemReq.Status == "Accept" {
		txn, err := db.Begin()
		if err != nil {
			return nil, err
		}
		defer txn.Rollback()

		res, err := txn.Exec("UPDATE users SET coins=coins-($1) WHERE rollno=($2) AND coins>=($1);", redeemReq.Coins, redeemReq.User)
		if err != nil {
			return nil, fmt.Errorf("Could not accept redeem request for user#%v; %v", redeemReq.User, err)
		} else if cntRows, err := res.RowsAffected(); err != nil {
			return nil, fmt.Errorf("Could not accept redeem request for user#%v; %v", redeemReq.User, err)
		} else if cntRows != 1 {
			return nil, fmt.Errorf("Could not accept redeem request for user#%v; Possible errors - insufficient funds in redeemer's account or user may have been deleted from the database", redeemReq.User)
		}

		res, err = txn.Exec("UPDATE redeemRequests SET status=1, description=(?), responded_on=CURRENT_TIMESTAMP WHERE id=(?) AND status NOT IN (0, 1);", redeemReq.Descr, redeemReq.Id)
		if err != nil {
			return nil, fmt.Errorf("Could not accept redeem request for user#%v; %v", redeemReq.User, err)
		} else if cntRows, err := res.RowsAffected(); err != nil {
			return nil, fmt.Errorf("Could not accept redeem request for user#%v; %v", redeemReq.User, err)
		} else if cntRows != 1 {
			return nil, fmt.Errorf("Could not accept redeem request for user#%v", redeemReq.User)
		}

		err = txn.Commit()
		if err != nil {
			return nil, err
		}
		return redeemReq.Id, nil
	} else {
		res, err := db.Exec("UPDATE redeemRequests SET status=0, description=(?), responded_on=CURRENT_TIMESTAMP WHERE id=(?) AND status NOT IN (0, 1);", redeemReq.Descr, redeemReq.Id)
		if err != nil {
			return nil, fmt.Errorf("Could not reject redeem request for user#%v; %v", redeemReq.User, err)
		} else if cntRows, err := res.RowsAffected(); err != nil {
			return nil, fmt.Errorf("Could not reject redeem request for user#%v; %v", redeemReq.User, err)
		} else if cntRows != 1 {
			return nil, fmt.Errorf("Could not reject redeem request for user#%v; Possible error - status may have already been updated", redeemReq.User)
		}
		return nil, nil
	}
}
