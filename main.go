package main

import (
	"log"
	"net/http"

	"github.com/AsishMandoi/iitk-coin/handlers"
	"github.com/AsishMandoi/iitk-coin/server"
)

func main() {

	go server.Routine("signup", handlers.Signup)

	go server.Routine("login", handlers.Login)

	go server.Routine("secret_page", handlers.Secret)

	go server.Routine("view_coins", handlers.ViewCoins)

	go server.Routine("transact", handlers.TransactCoins)

	go server.Routine("reward_coins", handlers.RewardCoins)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
