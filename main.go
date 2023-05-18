package main

import (
	"net/http"
	"os"
	"wine-project-go/repositories/categories"
	"wine-project-go/repositories/dishes"
	"wine-project-go/repositories/drinks"
	"wine-project-go/repositories/restaurants"
	"wine-project-go/repositories/users"

	"wine-project-go/repositories/login"
	"wine-project-go/repositories/register"

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

	router.HandleFunc("/login", login.Login).Methods("POST")

	router.HandleFunc("/users", register.CreateUserHandler).Methods("POST")
	router.HandleFunc("/users", users.GetUsers).Methods("GET")
	router.HandleFunc("/users/{id}", users.DeleteUserHandler).Methods("DELETE")
	router.HandleFunc("/users/{id}", users.UpdateUserHandler).Methods("PUT")

	router.HandleFunc("/category", categories.GetCategories).Methods("GET")
	router.HandleFunc("/category", categories.AddCategory).Methods("POST")
	//router.HandleFunc("/category", categories.AddCategory).Methods("DELETE")
	//router.HandleFunc("/category", categories.AddCategory).Methods("PUT")

	router.HandleFunc("/restaurant", restaurants.AddRestaurant).Methods("POST")
	router.HandleFunc("/restaurant", restaurants.GetRestaurants).Methods("GET")
	//router.HandleFunc("/category", categories.AddCategory).Methods("DELETE")
	//router.HandleFunc("/category", categories.AddCategory).Methods("PUT")

	router.HandleFunc("/dishes", dishes.AddDish).Methods("POST")
	router.HandleFunc("/dishes", dishes.GetDishes).Methods("GET")
	router.HandleFunc("/drinks", drinks.AddDrink).Methods("POST")
	router.HandleFunc("/drinks", drinks.GetDrinks).Methods("GET")

	handler := handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(router)
	http.ListenAndServe(":"+port, handler)
}
