package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/AsishMandoi/iitk-coin/database"
	"github.com/AsishMandoi/iitk-coin/global"
	"github.com/AsishMandoi/iitk-coin/server"
)

// GET request format
// --header 'Authorization: Bearer qWd3EjkVn&e6n.kJfvm82s3Fo@~389r$dml3%0v.s*Hsi&2-Y4'
func ViewCoins(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	payload := &global.ViewCoinsRespBody{} // Body of the response to be sent

	// Only accepting GET request from this endpoint
	if r.Method == "GET" {

		// Authorizing the request and obtaining the user's roll no
		statusCode, claims, err := server.ValidateJWT(r)
		if err != nil {
			server.Respond(w, payload, statusCode, nil, err.Error(), nil)
			return
		}

		// User rollno present in the token
		userRollno := int(claims["rollno"].(float64))

		// Handle initialization errors in SQLite DB
		if msg, err := database.InitMsg, database.InitErr; err != nil {
			server.Respond(w, payload, 500, msg, err.Error(), nil)
			return
		}

		coins, err := database.GetCoins(userRollno)
		if err == sql.ErrNoRows {
			server.Respond(w, payload, 400, fmt.Sprintf("Could not identify user with roll no %v", userRollno), err.Error(), nil)
			return
		} else if err != nil {
			server.Respond(w, payload, 500, "Could not fetch coins for the user", err.Error(), nil)
			return
		}
		server.Respond(w, payload, 200, "SUCCESS", nil, coins)
	} else {
		server.Respond(w, payload, 501, "Welcome to /view_coins! Please use a GET method to check your balance.", nil, nil)
	}
}
