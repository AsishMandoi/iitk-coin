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

// GET request format (in the body) -> {"rollno": 190197}
func ViewCoins(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	payload := &global.ViewCoinsRespBodyFormat{} // Body of the response to be sent

	if r.Method == "GET" {

		body := struct {
			Rollno int
		}{}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			server.Respond(w, payload, 400, "Could not decode body of the request", err.Error(), "-")
			return
		}

		// Initialize DB
		if msg, err := database.Initialize(); err != nil {
			server.Respond(w, payload, 500, msg, err.Error(), "-")
			return
		}

		coins, err := database.GetCoins(body.Rollno)
		if err != nil {
			if err == sql.ErrNoRows {
				server.Respond(w, payload, 400, fmt.Sprintf("Could not identify user with given roll no %v", body.Rollno), err.Error(), "-")
				return
			}
			server.Respond(w, payload, 500, "Something went wrong :(", err.Error(), "-")
			return
		}
		server.Respond(w, payload, 200, "SUCCESS", "-", coins)
	} else {
		server.Respond(w, payload, 501, "Welcome to /view_coins! Please use a GET request to check your balance.", "-", "-")
	}
}
