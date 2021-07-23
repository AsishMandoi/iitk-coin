package server

import (
	"crypto/rand"
	"crypto/tls"
	"math/big"
	"strconv"

	"github.com/AsishMandoi/iitk-coin/database"
	"github.com/AsishMandoi/iitk-coin/global"
	"gopkg.in/mail.v2"
)

func genOTP() (string, error) {
	otpNum, err := rand.Int(rand.Reader, big.NewInt(899999))
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(otpNum.Int64()+100000, 10), nil
}

// SendOTP performs these functions
// 1) generate OTP, 2) save it along with other details in the database with an expiry time and, 3) the send it
func SendOTP(emailid string, details interface{}, job string) (string, error) {

	otpStr, err := genOTP()
	if err != nil {
		return "Transaction failed: Could not generate OTP", err
	}

	var tfr global.TxnObj
	var rdm global.RedeemObj

	// STORE ALL DETAILS along with the otp IN THE REDIS DB with an expiry time
	if job == "transfer" {
		tfr = details.(global.TxnObj)
		tfr.Otp = otpStr
		if err = database.SetTfrDetails(tfr); err != nil {
			return "Transaction failed: Could not store OTP", err
		}
	} else {
		rdm = details.(global.RedeemObj)
		rdm.Otp = otpStr
		if err = database.SetRdmDetails(rdm); err != nil {
			return "Transaction failed: Could not store OTP", err
		}
	}

	gmailHost, gmailPort := "smtp.gmail.com", 587
	// iitkHost, iitkPort := "mmtp.iitk.ac.in", 25

	m := mail.NewMessage()
	m.SetHeader("From", global.MyGmailId)
	m.SetHeader("To", emailid)
	if job == "transfer" {
		m.SetHeader("Subject", "OTP for Transfer")
	} else {
		m.SetHeader("Subject", "OTP for Redeem")
	}

	// Plain text body
	// m.SetBody("text/plain", "This is a test email for IITK-Coin. There is no actual transaction going on.\nThe OTP for your transaction is\n"+otpStr+"\nThis OTP will expire in 2 minutes.\nDO NOT share it with anyone.\n\nIf this OTP was not requested by you, make sure to change your password immediately.\n\n\nRegards,\nAsish\nIITK-Coin")

	// HTML body
	m.SetBody("text/html", "<p><i>This is a test email for IITK-Coin. There is no actual transaction going on.</i></p><p>The OTP for your transaction is<br><h2>"+otpStr+"</h2></p><p>This OTP will expire in 2 minutes. DO NOT share it with anyone.<br><br>If this OTP was not requested by you, make sure to <b>change your password immediately</b>.<br><br><br>Regards,<br>Asish<br>IITK-Coin</p>")

	d := mail.NewDialer(gmailHost, gmailPort, global.MyGmailId, global.MyPwd)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		return "Transaction failed: Could not send OTP", err
	}

	return "", nil
}
