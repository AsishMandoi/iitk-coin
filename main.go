package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AsishMandoi/iitk-coin/handlers"
)

func main() {

	http.HandleFunc("/signup", handlers.Signup)

	http.HandleFunc("/login", handlers.Login)

	http.HandleFunc("/secret_page", handlers.Secret)

	http.HandleFunc("/view_coins", handlers.ViewCoins)

	http.HandleFunc("/transfer_coins", handlers.TransferCoins)

	http.HandleFunc("/reward_coins", handlers.RewardCoins)

	// Redeem related endpts
	http.HandleFunc("/redeem", handlers.Redeem)

	http.HandleFunc("/redeem_requests", handlers.ViewRedeemRequests)

	http.HandleFunc("/update_redeem_status", handlers.UpdateRedeemStatus)

	http.HandleFunc("/redeem_status", handlers.ViewRedeemStatus)

	fmt.Println("Server running on localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
