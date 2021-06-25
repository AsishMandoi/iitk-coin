package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AsishMandoi/iitk-coin/database"
	"github.com/AsishMandoi/iitk-coin/global"
	"github.com/AsishMandoi/iitk-coin/server"
)

// POST request format (in the body) -> {"sender": 190184, "receiver": 190197, "amount": 500}
func TransferCoins(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	payload := &global.DefaultRespBodyFormat{} // Body of the response to be sent

	if r.Method == "POST" {

		body := struct {
			Receiver int     `json:"receiver"`
			Amount   float64 `json:"amount"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			server.Respond(w, payload, 400, "Could not decode body of the request", err.Error())
			return
		}

		// Authorizing the request
		statusCode, claims, err := server.Authorize(r)
		if err != nil {
			server.Respond(w, payload, statusCode, nil, err.Error(), nil)
			return
		}

		// Initialize DB
		if msg, err := database.Initialize(); err != nil {
			server.Respond(w, payload, 500, msg, err.Error())
			return
		}

		sender := int(claims["rollno"].(float64))
		senderBatch := claims["batch"].(string)

		// Implement logic for the minimum events criteria for sender. Also store the variable MIN_EVENTS in .env
		// Implement logic for the minimum events criteria for sender. Also store the variable MIN_EVENTS in .env
		// Implement logic for the minimum events criteria for sender. Also store the variable MIN_EVENTS in .env
		// Implement logic for the minimum events criteria for sender. Also store the variable MIN_EVENTS in .env
		// Implement logic for the minimum events criteria for sender. Also store the variable MIN_EVENTS in .env

		receiverBatch, receiverRole, err := database.GetBatchnRole(body.Receiver)
		if err != nil {
			server.Respond(w, payload, 400, fmt.Sprintf("Transaction failed; Could not identify receiver with roll no %v", body.Receiver), err.Error())
			return
		}
		if receiverRole == "Admin" || receiverRole == "CoreTeam" {
			server.Respond(w, payload, 400, "Transaction failed", "Cannot tranfer coins to Admin or Core Team member")
		}

		recAmt := body.Amount * 0.98
		// How much coins are to be transferred depends on the batch
		if senderBatch != receiverBatch {
			recAmt = body.Amount * 0.67
		}

		if err := database.Transact(struct {
			Sender   int
			Receiver int
			Amount   float64
		}{sender, body.Receiver, body.Amount}, recAmt); err != nil {
			server.Respond(w, payload, 400, "Transaction failed", err.Error())
			return
		}
		server.Respond(w, payload, 200, fmt.Sprintf("Transaction Successful; User: #%v transferred %v coins to user: #%v", sender, recAmt, body.Receiver), nil)
	} else {
		server.Respond(w, payload, 501, "Welcome to /transfer_coins page! Please use a POST request to send coins to another user.", nil)
	}
}
