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

	fmt.Println("Server running on localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
