package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AsishMandoi/iitk-coin/database"
	"github.com/AsishMandoi/iitk-coin/global"
	"github.com/AsishMandoi/iitk-coin/server"
)

func RewardCoins(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	payload := &global.TxnRespBody{} // Body of the response to be sent

	if r.Method == "POST" {

		body := struct {
			Receiver int     `json:"receiver"`
			Amount   float64 `json:"amount"`
			Descr    string  `json:"description"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			server.Respond(w, payload, 400, "Could not decode body of the request", err.Error(), nil)
			return
		}

		// Authorizing the request and obtaining the user's roll no
		statusCode, claims, err := server.ValidateJWT(r)
		if err != nil {
			server.Respond(w, payload, statusCode, nil, err.Error(), nil)
			return
		}

		// Making sure that an Admin is rewarding coins
		if claims["role"] != "Admin" {
			server.Respond(w, payload, 401, nil, "User unauthorized to reward coins", nil)
			return
		}

		sender := int(claims["rollno"].(float64))

		// Handle initialization errors in SQLite DB
		if msg, err := database.InitMsg, database.InitErr; err != nil {
			server.Respond(w, payload, 500, msg, err.Error(), nil)
			return
		}

		// Making sure that the coins are being rewarded to a valid user but not to an Admin
		if _, receiverRole, err := database.GetBatchnRole(body.Receiver); err != nil {
			if err == sql.ErrNoRows {
				server.Respond(w, payload, 400, fmt.Sprintf("Reward Failed; Could not identify receiver with roll no %v", body.Receiver), err.Error(), nil)
			} else {
				server.Respond(w, payload, 400, "Reward Failed", err.Error(), nil)
			}
			return
		} else if receiverRole == "Admin" {
			server.Respond(w, payload, 400, "Reward Failed", "Cannot reward coins to an Admin", nil)
			return
		}

		if txnId, err := database.Reward(global.TxnObj{sender, body.Receiver, body.Amount, body.Amount, body.Descr, ""}); err != nil {
			server.Respond(w, payload, 400, "Reward failed", err.Error(), nil)
		} else {
			server.Respond(w, payload, 200, fmt.Sprintf("Reward Successful; User: #%v was rewarded with %v coins", body.Receiver, body.Amount), nil, txnId)
		}
	} else {
		server.Respond(w, payload, 501, "Welcome to /reward page! Please use a POST method to reward a user.", nil, nil)
	}
}
