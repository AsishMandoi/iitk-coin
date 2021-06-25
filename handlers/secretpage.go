package handlers

import (
	"net/http"

	"github.com/AsishMandoi/iitk-coin/global"
	"github.com/AsishMandoi/iitk-coin/server"
)

// GET request format (in the header) -> --header "Authorization: Bearer <access token>"
func Secret(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	payload := &global.SecretpageRespBodyFormat{} // Body of the response to be sent

	if r.Method == "GET" {

		// Authorizing the request
		statusCode, claims, err := server.Authorize(r)
		if err != nil {
			server.Respond(w, payload, statusCode, "User unauthorized", err.Error(), nil)
			return
		}

		// Since there are no more errors, the secretpage responds with the confidential information.
		server.Respond(w, payload, statusCode, "SUCCESS", nil, int(claims["rollno"].(float64)))
	} else {
		server.Respond(w, payload, 501, "Welcome to /secret_page! Please use a GET request to get authorized.", nil, nil)
	}
}
