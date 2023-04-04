package main

import (
	"db-go/login/register"
	"log"
	"net/http"
)

func main() {

	srv := http.Server{
		Addr: ":8081",
	}

	http.HandleFunc("/signup", register.CreateUserHandler)

	log.Fatal(srv.ListenAndServe())
}
