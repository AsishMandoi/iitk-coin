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
	payload := &global.DefaultRespBody{} // Body of the response to be sent

	if r.Method == "POST" {

		body := struct {
			Receiver int     `json:"receiver"`
			Amount   float64 `json:"amount"`
			Descr    string  `json:"description"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			server.Respond(w, payload, 400, "Could not decode body of the request", err.Error())
			return
		}

		// Authorizing the request and obtaining the user's roll no
		statusCode, claims, err := server.ValidateJWT(r)
		if err != nil {
			server.Respond(w, payload, statusCode, nil, err.Error())
			return
		}

		// Handle initialization errors in DB
		if msg, err := database.InitMsg, database.InitErr; err != nil {
			server.Respond(w, payload, 500, msg, err.Error())
			return
		}

		senderRoll := int(claims["rollno"].(float64))
		senderBatch := claims["batch"].(string)
		// fmt.Printf("%v->(%T)\n\n\n", claims["email"], claims["email"])
		senderEmail := claims["email"].(string)

		receiverBatch, receiverRole, err := database.GetBatchnRole(body.Receiver)
		if err != nil {
			server.Respond(w, payload, 400, fmt.Sprintf("Transaction failed: could not identify receiver with roll no %v", body.Receiver), err.Error())
			return
		}
		if receiverRole == "Admin" || receiverRole == "CoreTeam" {
			server.Respond(w, payload, 400, "Transaction failed", "Cannot tranfer coins to an Admin or a Core Team member")
			return
		}

		// The sender and the receiver both need to have participated in at least `MinEvents` number of events.
		if cntEvents, err := database.GetCntEvents(senderRoll); err != nil {
			server.Respond(w, payload, 400, "Transaction failed", err.Error())
			return
		} else if cntEvents < global.MinEvents {
			server.Respond(w, payload, 403, "Transaction failed: sender not eligible", fmt.Sprintf("User #%v is not yet eligible to make transactions", senderRoll))
			return
		}

		if cntEvents, err := database.GetCntEvents(body.Receiver); err != nil {
			server.Respond(w, payload, 400, "Transaction failed", err.Error())
			return
		} else if cntEvents < global.MinEvents {
			server.Respond(w, payload, 403, "Transaction failed: receiver not eligible", fmt.Sprintf("User #%v is not yet eligible to participate in transactions", body.Receiver))
			return
		}

		amtRcvd := body.Amount * 0.98
		// How much coins are to be transferred depends on the sender's as well as the receiver's batch
		if senderBatch != receiverBatch {
			amtRcvd = body.Amount * 0.67
		}

		// Generate OTP, save it (along with other details) in the database with an expiry time and then send it
		if msg, err := server.SendOTP(senderEmail, global.TxnObj{senderRoll, body.Receiver, body.Amount, amtRcvd, body.Descr, ""}, "transfer"); err != nil {
			server.Respond(w, payload, 500, msg, err.Error())
			return
		}

		server.Respond(w, payload, 200, "Post your otp on http://localhost:8080/confirm_transfer to confirm your transaction", nil)

		// REDIRECT WILL NOT WORK
		// http.Redirect(w, r, "/confirm_transfer", http.StatusPermanentRedirect)

		// ???? http.RedirectHandler()
	} else {
		server.Respond(w, payload, 501, "Welcome to /transfer_coins page! Please use a POST method to send coins to another user.", nil)
	}
}
