package main

import (
	"net/http"
	"os"
	"wine-project-go/categories"

	"wine-project-go/login"
	"wine-project-go/login/register"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8081"
	}

	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})

	router := mux.NewRouter()

	router.HandleFunc("/signup", register.CreateUserHandler).Methods("POST")
	router.HandleFunc("/login", login.Login).Methods("POST")
	router.HandleFunc("/category", categories.AddCategory).Methods("POST")
	router.HandleFunc("/category", categories.GetCategories).Methods("GET")

	// authRouter := router.PathPrefix("/api").Subrouter()
	// authRouter.Use(session.AuthMiddleware)
	// authRouter.HandleFunc("/users", ListUsersHandler).Methods("GET")

	handler := handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(router)
	http.ListenAndServe(":"+port, handler)
}
