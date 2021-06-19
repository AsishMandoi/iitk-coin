package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AsishMandoi/iitk-coin/database"
	"github.com/AsishMandoi/iitk-coin/global"
	"github.com/AsishMandoi/iitk-coin/server"
)

// POST request format (in the body) -> {"receiver": 190197, "amount": 500}
func RewardCoins(w http.ResponseWriter, r *http.Request) {
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

		// Initialize DB
		if msg, err := database.Initialize(); err != nil {
			server.Respond(w, payload, 500, msg, err.Error())
			return
		}

		err := database.Reward(body)
		if err != nil {
			server.Respond(w, payload, 400, "Reward failed", err.Error())
			return
		}
		server.Respond(w, payload, 200, fmt.Sprintf("Reward Successful; User: #%v was rewarded with %v coins", body.Receiver, body.Amount), "-")
	} else {
		server.Respond(w, payload, 501, "Welcome to /reward_coins page! Please use a POST request to reward a user.", "-")
	}
}
