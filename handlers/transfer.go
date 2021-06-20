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
			Sender   int     `json:"sender"`
			Receiver int     `json:"receiver"`
			Amount   float64 `json:"amount"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			server.Respond(w, payload, 400, "Could not decode body of the request", err.Error())
			return
		}

		// Initialize DB
		if msg, err := database.Initialize(); err != nil {
			server.Respond(w, payload, 500, msg, err.Error())
			return
		}

		// Logic to calculate coins to be transferred depending on the batch
		senderBatch, err := database.GetBatch(body.Sender)
		if err != nil {
			server.Respond(w, payload, 400, "Transaction failed", fmt.Errorf("Could not send amount; %v", err))
			return
		}
		receiverBatch, err := database.GetBatch(body.Receiver)
		if err != nil {
			server.Respond(w, payload, 400, "Transaction failed", fmt.Errorf("Could not receive amount; %v", err))
			return
		}
		recAmt := body.Amount * 0.98
		if senderBatch != receiverBatch {
			recAmt = body.Amount * 0.67
		}

		if err := database.Transact(body, recAmt); err != nil {
			server.Respond(w, payload, 400, "Transaction failed", err.Error())
			return
		}
		server.Respond(w, payload, 200, fmt.Sprintf("Transaction Successful; User: #%v transferred %v coins to user: #%v", body.Sender, recAmt, body.Receiver), "-")
	} else {
		server.Respond(w, payload, 501, "Welcome to /transact page! Please use a POST request to send coins to another user.", "-")
	}
}
