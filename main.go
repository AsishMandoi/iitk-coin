package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AsishMandoi/iitk-coin/handlers"
)

func main() {
	http.HandleFunc("/signup", handlers.Signup)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/reset_password", handlers.ResetPwd)
	http.HandleFunc("/secret_page", handlers.Secret)
	http.HandleFunc("/view_coins", handlers.ViewCoins)
	http.HandleFunc("/transfer", handlers.TransferCoins)
	http.HandleFunc("/reward", handlers.RewardCoins)

	// Redeem related endpoints
	http.HandleFunc("/redeems/new", handlers.NewRedeemReq)
	http.HandleFunc("/redeems", handlers.ViewRedeemRequests)
	http.HandleFunc("/redeems/update", handlers.UpdateRedeemStatus)
	http.HandleFunc("/redeems/status", handlers.ViewRedeemStatus)

	// Confirmation endpoints
	http.HandleFunc("/reset_password/confirm", handlers.ResetPwdCheck)
	http.HandleFunc("/transfer/confirm", handlers.TransferCheck)
	http.HandleFunc("/redeems/new/confirm", handlers.RedeemCheck)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server running on port:" + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
