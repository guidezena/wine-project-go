package main

import (
	"db-go/login"
	"db-go/login/register"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/signup", register.CreateUserHandler)
	router.HandleFunc("/login", login.Login)

	// authRouter := router.PathPrefix("/api").Subrouter()
	// authRouter.Use(session.AuthMiddleware)
	// authRouter.HandleFunc("/users", ListUsersHandler).Methods("GET")

	http.ListenAndServe(":8081", router)
}
