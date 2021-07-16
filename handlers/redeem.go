package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AsishMandoi/iitk-coin/database"
	"github.com/AsishMandoi/iitk-coin/global"
	"github.com/AsishMandoi/iitk-coin/server"
)

// Making a redeem request
func Redeem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	payload := &global.RedeemRespBody{} // Body of the response to be sent
	if r.Method == "POST" {
		body := struct {
			Item_id int     `json:"item_id"`
			Price   float64 `json:"price"`
			Descr   string  `json:"description"`
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

		usr := int(claims["rollno"].(float64))

		// Handle initialization errors in DB
		if msg, err := database.InitMsg, database.InitErr; err != nil {
			server.Respond(w, payload, 500, msg, err.Error(), nil)
			return
		}

		if reqId, err := database.RedeemReq(usr, body); err != nil {
			server.Respond(w, payload, 400, "Redeem request failed", err.Error(), nil)
		} else {
			server.Respond(w, payload, 201, "Redeem request successful", nil, reqId)
		}
	} else {
		server.Respond(w, payload, 501, "Welcome to /redeem page! Please use a POST method to make a redeem request.", nil, nil)
	}
}

// Update the status of redeem requests (Admin only)
func UpdateRedeemStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	payload := &global.TxnRespBody{} // Body of the response to be sent
	if r.Method == "POST" {
		body := global.RedeemStatusUPDBody{}
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

		// Making sure that an Admin is updating the redeem request status
		if claims["role"] != "Admin" {
			server.Respond(w, payload, 401, nil, "User unauthorized to update redeem status", nil)
			return
		}

		// Handle initialization errors in DB
		if msg, err := database.InitMsg, database.InitErr; err != nil {
			server.Respond(w, payload, 500, msg, err.Error(), nil)
			return
		}

		if txnId, err := database.UpdRdmSts(body); err != nil {
			server.Respond(w, payload, 400, "Redeem status updation failed", err.Error(), nil)
		} else {
			server.Respond(w, payload, 200, "Redeem updated successfully", nil, txnId)
		}

	} else {
		server.Respond(w, payload, 501, "Welcome to /update_redeem_status page! Please use a POST method to update the status of a redeem request.", nil, nil)
	}
}

// Check the status of all redeem requests
// Will later add a functionality of - show the status of requests before/after a certain date or between 2 dates
func ViewRedeemStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	payload := &global.DefaultDataRespBody{} // Body of the response to be sent
	if r.Method == "GET" {

		// Authorizing the request and obtaining the user's roll no
		statusCode, claims, err := server.ValidateJWT(r)
		if err != nil {
			server.Respond(w, payload, statusCode, nil, err.Error(), nil)
			return
		}

		usr := int(claims["rollno"].(float64))

		// Handle initialization errors in DB
		if msg, err := database.InitMsg, database.InitErr; err != nil {
			server.Respond(w, payload, 500, msg, err.Error(), nil)
			return
		}

		if data, err := database.ShowRdmSts(usr); err != nil {
			server.Respond(w, payload, 400, "Could not fetch redeem info", err.Error(), nil)
		} else {
			server.Respond(w, payload, 200, nil, nil, data)
		}

	} else {
		server.Respond(w, payload, 501, "Welcome to /redeem_status page! Please use a GET method to view the status of a redeem request.", nil, nil)
	}
}

// See all pending redeem requests (Admin only)
// Will later add a functionality of - show the pending requests before/after a certain date or between 2 dates
func ViewRedeemRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	payload := &global.DefaultDataRespBody{} // Body of the response to be sent
	if r.Method == "GET" {

		// Authorizing the request and obtaining the user's roll no
		statusCode, claims, err := server.ValidateJWT(r)
		if err != nil {
			server.Respond(w, payload, statusCode, nil, err.Error(), nil)
			return
		}

		// Making sure that an Admin is updating the redeem request status
		if claims["role"] != "Admin" {
			server.Respond(w, payload, 401, nil, "User unauthorized to view redeem requests", nil)
			return
		}

		// Handle initialization errors in DB
		if msg, err := database.InitMsg, database.InitErr; err != nil {
			server.Respond(w, payload, 500, msg, err.Error(), nil)
			return
		}

		if data, err := database.ShowAllRdmReqs(); err != nil {
			server.Respond(w, payload, 400, "Cannot show redeem requests", err.Error(), nil)
		} else {
			server.Respond(w, payload, 200, nil, nil, data)
		}

	} else {
		server.Respond(w, payload, 501, "Welcome to /redeem_requests page! Please use a GET method to view pending redeem requests.", nil, nil)
	}
}
