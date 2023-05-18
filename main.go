package main

import (
	"net/http"
	"os"
	"wine-project-go/repositories/categories"
	"wine-project-go/repositories/dishes"
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

	router.HandleFunc("/categories", categories.GetCategoriesHandler).Methods("GET")
	router.HandleFunc("/categories", categories.AddCategoryHandler).Methods("POST")
	router.HandleFunc("/categories/{id}", categories.DeleteCategoryHandler).Methods("DELETE")
	router.HandleFunc("/categories/{id}", categories.UpdateCategoryHandler).Methods("PUT")

	router.HandleFunc("/restaurants", restaurants.AddRestaurantHandler).Methods("POST")
	router.HandleFunc("/restaurants", restaurants.GetRestaurantsHandler).Methods("GET")
	router.HandleFunc("/restaurants/{id}", restaurants.DeleteRestaurantHandler).Methods("DELETE")
	router.HandleFunc("/restaurants/{id}", restaurants.UpdateRestaurantHandler).Methods("PUT")

	router.HandleFunc("/dishes", dishes.AddDishHandler).Methods("POST")
	router.HandleFunc("/dishes", dishes.GetDishesHandler).Methods("GET")
	router.HandleFunc("/dishes/{id}", dishes.DeleteDishHandler).Methods("DELETE")
	router.HandleFunc("/dishes/{id}", dishes.UpdateDishHandler).Methods("PUT")

	//router.HandleFunc("/drinks", drinks.AddDrinkHandler).Methods("POST")
	//router.HandleFunc("/drinks", drinks.GetDrinksHandler).Methods("GET")
	//router.HandleFunc("/drinks/{id}", drinks.DeleteDrinkHandler).Methods("DELETE")
	//router.HandleFunc("/drinks/{id}", drinks.UpdateDrinksHandler).Methods("PUT")

	handler := handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(router)
	http.ListenAndServe(":"+port, handler)
}
