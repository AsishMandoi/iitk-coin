package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AsishMandoi/iitk-coin/database"
	"github.com/AsishMandoi/iitk-coin/global"
	"github.com/AsishMandoi/iitk-coin/server"
)

// POST request format
// --header 'Authorization: Bearer qWd3EjkVn-e6n.kJfvm82s3Fo@~389r$dml3-0v.s*Hsi&2-Y4'
// --data '{"receiver": 190197, "amount": 500, "desription": "samose kha lena"}'
func TransferCoins(w http.ResponseWriter, r *http.Request) {
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
			server.Respond(w, payload, statusCode, nil, err.Error(), nil, nil)
			return
		}

		// Handle initialization errors in DB
		if msg, err := database.InitMsg, database.InitErr; err != nil {
			server.Respond(w, payload, 500, msg, err.Error(), nil)
			return
		}

		sender := int(claims["rollno"].(float64))
		senderBatch := claims["batch"].(string)

		receiverBatch, receiverRole, err := database.GetBatchnRole(body.Receiver)
		if err != nil {
			server.Respond(w, payload, 400, fmt.Sprintf("Transaction failed; Could not identify receiver with roll no %v", body.Receiver), err.Error(), nil)
			return
		}
		if receiverRole == "Admin" || receiverRole == "CoreTeam" {
			server.Respond(w, payload, 400, "Transaction failed", "Cannot tranfer coins to an Admin or a Core Team member", nil)
			return
		}

		// The sender and the receiver both need to have participated in at least `MinEvents` number of events.
		if cntEvents, err := database.GetCntEvents(sender); err != nil {
			server.Respond(w, payload, 400, "Transaction failed", err.Error(), nil)
			return
		} else if cntEvents < global.MinEvents {
			server.Respond(w, payload, 403, "Transaction failed; Sender not eligible", fmt.Sprintf("User #%v is not yet eligible to make transactions", sender), nil)
			return
		}

		if cntEvents, err := database.GetCntEvents(body.Receiver); err != nil {
			server.Respond(w, payload, 400, "Transaction failed", err.Error(), nil)
			return
		} else if cntEvents < global.MinEvents {
			server.Respond(w, payload, 403, "Transaction failed; Receiver not eligible", fmt.Sprintf("User #%v is not yet eligible to participate in transactions", body.Receiver), nil)
			return
		}

		amtRcvd := body.Amount * 0.98
		// How much coins are to be transferred depends on the sender's as well as the receiver's batch
		if senderBatch != receiverBatch {
			amtRcvd = body.Amount * 0.67
		}

		if txnId, err := database.Transact(global.TxnBody{sender, body.Receiver, body.Amount, body.Descr}, amtRcvd); err != nil {
			server.Respond(w, payload, 400, "Transaction failed", err.Error(), nil)
		} else {
			server.Respond(w, payload, 200, fmt.Sprintf("Transaction Successful; User: #%v transferred %v coins to user: #%v", sender, amtRcvd, body.Receiver), nil, txnId)
		}
	} else {
		server.Respond(w, payload, 501, "Welcome to /transfer_coins page! Please use a POST method to send coins to another user.", nil, nil)
	}
}
