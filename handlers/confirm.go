package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AsishMandoi/iitk-coin/database"
	"github.com/AsishMandoi/iitk-coin/global"
	"github.com/AsishMandoi/iitk-coin/server"
)

func TransferCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	payload := &global.TxnRespBody{} // Body of the response to be sent

	if r.Method == "POST" {

		body := struct {
			Otp       string `json:"otp"`
			ResendOTP bool   `json:"resend"`
		}{}

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

		sender := int(claims["rollno"].(float64))

		// Transfer details are collected from the database
		tfrDet, msg, err := database.GetTfrDetails(sender)
		if err != nil {
			server.Respond(w, payload, 500, "Transaction failed: "+msg, err.Error(), nil)
			return
		}

		if body.ResendOTP {
			email := claims["email"].(string)
			if msg, err := server.SendOTP(email, tfrDet, "transfer"); err != nil {
				server.Respond(w, payload, 500, msg, err.Error(), nil)
				return
			}
			server.Respond(w, payload, 200, "OTP sent", nil, nil)
			return
		}

		// Validating the OTP (i.e. checking if the entered otp is same as that stored in database)
		if tfrDet.Otp != body.Otp {
			server.Respond(w, payload, 401, "Transaction failed", "Incorrect OTP", nil)
			return
		}

		// Handle initialization errors in DB
		if msg, err := database.InitMsg, database.InitErr; err != nil {
			server.Respond(w, payload, 500, msg, err.Error(), nil)
			return
		}

		if txnId, err := database.Transact(tfrDet); err != nil {
			server.Respond(w, payload, 500, "Transaction failed", err.Error(), nil)
		} else {
			server.Respond(w, payload, 200, fmt.Sprintf("Transaction Successful: User: #%v transferred %v coins to user: #%v", sender, tfrDet.AmtRcvd, tfrDet.Receiver), nil, txnId)
			database.DelTfrDetails(sender)
		}
	} else {
		server.Respond(w, payload, 405, "Method not allowed for this endpoint", nil, nil)
	}
}

func RedeemCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	payload := &global.RedeemRespBody{} // Body of the response to be sent

	if r.Method == "POST" {

		body := struct {
			Otp       string `json:"otp"`
			ResendOTP bool   `json:"resend"`
		}{}

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

		redeemer := int(claims["rollno"].(float64))

		// Redeem details are collected from the database
		rdmDet, msg, err := database.GetRdmDetails(redeemer)
		if err != nil {
			server.Respond(w, payload, 500, "Redeem request failed: "+msg, err.Error(), nil)
			return
		}

		if body.ResendOTP {
			email := claims["email"].(string)
			if msg, err := server.SendOTP(email, rdmDet, "redeem"); err != nil {
				server.Respond(w, payload, 500, msg, err.Error(), nil)
				return
			}
			server.Respond(w, payload, 200, "OTP sent", nil, nil)
			return
		}

		// Validating the OTP (i.e. checking if the entered otp is same as that stored in database)
		if rdmDet.Otp != body.Otp {
			server.Respond(w, payload, 401, "Redeem request failed", "Incorrect OTP", nil)
			return
		}

		// Handle initialization errors in DB
		if msg, err := database.InitMsg, database.InitErr; err != nil {
			server.Respond(w, payload, 500, msg, err.Error(), nil)
			return
		}

		if reqId, err := database.RedeemReq(rdmDet); err != nil {
			server.Respond(w, payload, 500, "Redeem request failed", err.Error(), nil)
		} else {
			server.Respond(w, payload, 201, "Redeem request successful", nil, reqId)
			database.DelRdmDetails(redeemer)
		}
	} else {
		server.Respond(w, payload, 405, "Method not allowed for this endpoint", nil, nil)
	}
}
