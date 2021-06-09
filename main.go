package main

import (
	"log"
	"net/http"
)

type Stu struct {
	Rollno   int
	Name     string
	Password string
}

type secretpageRespBodyFormat struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Data    string `json:"data"`
}

type signupRespBodyFormat struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type loginInputFormat struct {
	Rollno   int    `json:"rollno"`
	Password string `json:"password"`
}

type loginRespBodyFormat struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Token   string `json:"token"`
}

func main() {

	http.HandleFunc("/signup", signup)

	http.HandleFunc("/login", login)

	http.HandleFunc("/secretpage", validateToken)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
