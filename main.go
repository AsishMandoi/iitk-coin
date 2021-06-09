package main

import (
	"log"
	"net/http"

	"github.com/AsishMandoi/iitk-coin/handlers"
)

func main() {

	http.HandleFunc("/signup", handlers.Signup)

	http.HandleFunc("/login", handlers.Login)

	http.HandleFunc("/secretpage", handlers.ValidateToken)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
