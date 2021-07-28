package database

import (
	"fmt"

	"github.com/AsishMandoi/iitk-coin/global"
	_ "github.com/mattn/go-sqlite3"
)

// Get the user details (password, batch, and role) for the given rollno (if present)
func GetUsrDetails(rollno int) (struct{ Email, Pwd, Batch, Role string }, error) {
	row := db.QueryRow("SELECT email, password, batch, role FROM users WHERE rollno=?;", rollno)
	var usrDetails struct{ Email, Pwd, Batch, Role string }
	err := row.Scan(&usrDetails.Email, &usrDetails.Pwd, &usrDetails.Batch, &usrDetails.Role)
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
	_, err := db.Exec("INSERT INTO users(rollno, name, email, password, batch, role, coins) VALUES (?, ?, ?, ?, ?, ?, ?);", usr.Rollno, usr.Name, usr.Email, usr.Password, usr.Batch, "", 0)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.rollno" {
			return "Could not add user into the database", fmt.Errorf("User #%v already present", usr.Rollno)
		}
		return "Could not add user into the database", err
	}
	return "Added user successfully", nil
}

func UpdPwd(rollno int, pwd string) error {
	res, err := db.Exec("UPDATE users SET password=(?) WHERE rollno=(?);", pwd, rollno)
	if err != nil {
		return err
	} else if cntRows, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("Could not update password for user#%v: %v", rollno, err)
	} else if cntRows != 1 {
		return fmt.Errorf("Could not update password for user#%v", rollno)
	}
	return nil
}

var reqId int
var usr int
var itemId int
var amt float64
var descr string
var rdmStatus = map[int]string{0: "Rejected", 1: "Accepted", 2: "Pending"}
var reqTime interface{}
var respTime interface{}

// Add an entry to the redeemRequests table in the DB
func RedeemReq(rdm global.RedeemObj) (interface{}, error) {
	res, err := db.Exec("INSERT INTO redeemRequests(redeemer, item_id, amount, description, status) VALUES (?, ?, ?, ?, 2);", rdm.Redeemer, rdm.ItemId, rdm.Amount, rdm.Descr)
	if err != nil {
		return nil, err
	} else if cntRows, err := res.RowsAffected(); err != nil {
		return nil, err
	} else if cntRows != 1 {
		return nil, fmt.Errorf("Could not make redeem request")
	}

	reqId, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("Could not make redeem request; %v", err)
	}
	return reqId, nil
}

// Show all pending redeem requests (to the admin) [A list of objects of type global.RedeemReqObj is returned]
func ShowAllRdmReqs() (interface{}, error) {
	rows, err := db.Query("SELECT id, redeemer, item_id, amount, description, requested_on FROM redeemRequests WHERE status=2")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var allRedeemReqs []global.RedeemReqObj

	for rows.Next() {
		if err := rows.Scan(&reqId, &usr, &itemId, &amt, &descr, &reqTime); err != nil {
			return nil, err
		}
		allRedeemReqs = append(allRedeemReqs, global.RedeemReqObj{reqId, usr, itemId, amt, descr, reqTime})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return allRedeemReqs, nil
}

// Show the statuses of all redeem requests to a user [A list of objects of type global.UserRedeemState is returned]
func ShowRdmSts(redeemer int) (interface{}, error) {
	rows, err := db.Query("SELECT id, item_id, amount, description, status, requested_on, responded_on FROM redeemRequests WHERE redeemer=(?);", redeemer)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var myRedeems []global.UserRedeemState
	var idx int

	for rows.Next() {
		if err := rows.Scan(&reqId, &itemId, &amt, &descr, &idx, &reqTime, &respTime); err != nil {
			return nil, err
		}
		myRedeems = append(myRedeems, global.UserRedeemState{reqId, itemId, amt, descr, rdmStatus[idx], reqTime, respTime})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return myRedeems, nil
}
