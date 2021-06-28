package handlers

import (
	"net/http"

	"github.com/AsishMandoi/iitk-coin/global"
	"github.com/AsishMandoi/iitk-coin/server"
)

// GET request format
// --header 'Authorization: Bearer qWd3EjkVn-e6n.kJfvm82s3Fo@~389r$dml3-0v.s*Hsi&2-Y4'
func Secret(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	payload := &global.SecretpageRespBody{} // Body of the response to be sent

	if r.Method == "GET" {

		// Authorizing the request
		if statusCode, claims, err := server.ValidateJWT(r); err != nil {
			server.Respond(w, payload, statusCode, "User unauthorized", err.Error(), nil)
		} else {
			server.Respond(w, payload, statusCode, "SUCCESS", nil, int(claims["rollno"].(float64)))
		}
	} else {
		server.Respond(w, payload, 501, "Welcome to /secret_page! Please use a GET request to get authorized.", nil, nil)
	}
}
