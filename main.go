package main

import (
	"db-go/login"
	"db-go/login/register"
	"log"
	"net/http"
)

func main() {

	srv := http.Server{
		Addr: ":8081",
	}

	http.HandleFunc("/signup", register.CreateUserHandler)
	http.HandleFunc("/login", login.Login)

	log.Fatal(srv.ListenAndServe())
}
