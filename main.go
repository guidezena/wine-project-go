package main

import (
	"net/http"
	"os"
	"wine-project-go/login"
	"wine-project-go/login/register"

	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	router := mux.NewRouter()

	router.HandleFunc("/signup", register.CreateUserHandler)
	router.HandleFunc("/login", login.Login)

	// authRouter := router.PathPrefix("/api").Subrouter()
	// authRouter.Use(session.AuthMiddleware)
	// authRouter.HandleFunc("/users", ListUsersHandler).Methods("GET")

	if port == "" {
		port = "8081"
	}
	http.ListenAndServe(":"+port, router)
}
