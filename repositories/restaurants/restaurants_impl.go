package restaurants

import (
	"encoding/json"
	"log"
	"net/http"
	"wine-project-go/dbConnection"
	"wine-project-go/entities"
	"wine-project-go/utils"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func AddRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request AddRestaurantHandler")

	var restaurant entities.Restaurant
	err := json.NewDecoder(r.Body).Decode(&restaurant)

	log.Printf(restaurant.Name)

	if err != nil {
		utils.SendError(w, "Decode error", http.StatusBadRequest)
		return
	}

	writer := dbConnection.GetWriterGorm()
	errorToWrite := createRestaurant(writer, restaurant)
	dbConnection.CloseDbConnection(writer)

	if errorToWrite != nil {
		log.Printf("errorToWrite")

		utils.SendError(w, errorToWrite.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]string{
		"message": "Restaurante criada com sucesso!",
	}

	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
	log.Printf("Restaurante criada com sucesso!")

}

func createRestaurant(db *gorm.DB, restaurant entities.Restaurant) error {
	log.Printf("createRestaurant")

	newRestaurant := entities.Restaurant{
		IdUser:      restaurant.IdUser,
		Name:        restaurant.Name,
		Image:       restaurant.Image,
		Description: restaurant.Description,
		Address:     restaurant.Address,
	}

	result := db.Create(&newRestaurant)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &entities.CustomError{Message: "Nenhuma linha foi afetada"}
	}

	return nil
}

func GetRestaurantsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request GetRestaurantsHandler")

	reader := dbConnection.GetReaderGorm()
	restaurants, err := getRestaurants(reader)
	dbConnection.CloseDbConnection(reader)

	if err != nil {
		http.Error(w, "Erro ao obter restaurantes", http.StatusInternalServerError)
		return
	}

	jsonRestaurants, err := json.Marshal(restaurants)
	if err != nil {
		http.Error(w, "Erro ao converter restaurantes para JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonRestaurants)
}

func getRestaurants(db *gorm.DB) ([]entities.Restaurant, error) {
	var restaurants []entities.Restaurant
	err := db.Find(&restaurants).Error
	if err != nil {
		return nil, err
	}
	return restaurants, nil
}

func GetRestaurantsForIdHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request GetRestaurantsHandler")

	restaurantID := mux.Vars(r)["id"]

	reader := dbConnection.GetReaderGorm()
	restaurants, err := getRestaurantsForId(reader, restaurantID)
	dbConnection.CloseDbConnection(reader)

	if err != nil {
		http.Error(w, "Erro ao obter restaurantes", http.StatusInternalServerError)
		return
	}

	jsonRestaurants, err := json.Marshal(restaurants)
	if err != nil {
		http.Error(w, "Erro ao converter restaurantes para JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonRestaurants)
}

func getRestaurantsForId(db *gorm.DB, restaurantId string) (*entities.Restaurant, error) {

	var restaurant entities.Restaurant
	result := db.First(&restaurant, restaurantId)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, &entities.CustomError{Message: "Nenhuma linha foi afetada"}
	}

	return &restaurant, nil

}

func DeleteRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request DeleteRestaurantHandler")

	restaurantID := mux.Vars(r)["id"]

	writer := dbConnection.GetWriterGorm()
	errorToDelete := deleteRestaurant(writer, restaurantID)
	dbConnection.CloseDbConnection(writer)

	if errorToDelete != nil {
		utils.SendError(w, errorToDelete.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"message": "Resturant excluída com sucesso!",
	}

	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	log.Printf("Restaurante excluído com sucesso!")
}

func deleteRestaurant(db *gorm.DB, restaurantID string) error {
	log.Printf("deleteRestaurant")

	result := db.Delete(&entities.Restaurant{}, restaurantID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &entities.CustomError{Message: "Nenhuma linha foi afetada"}
	}

	return nil
}

func UpdateRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request UpdateRestaurantHandler")

	restaurantID := mux.Vars(r)["id"]

	var updatedRestaurant entities.Restaurant
	err := json.NewDecoder(r.Body).Decode(&updatedRestaurant)

	if err != nil {
		utils.SendError(w, "Decode error", http.StatusBadRequest)
		return
	}

	writer := dbConnection.GetWriterGorm()
	errorToUpdate := updateRestaurant(writer, restaurantID, updatedRestaurant)
	dbConnection.CloseDbConnection(writer)

	if errorToUpdate != nil {
		utils.SendError(w, errorToUpdate.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"message": "Restaurante atualizado com sucesso!",
	}

	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	log.Printf("Restaurante atualizado com sucesso!")
}

func updateRestaurant(db *gorm.DB, restaurantID string, updatedRestaurant entities.Restaurant) error {
	log.Printf("updateRestaurant")

	result := db.Model(&entities.Restaurant{}).Where("ID = ?", restaurantID).Updates(updatedRestaurant)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &entities.CustomError{Message: "Nenhuma linha foi afetada"}
	}

	return nil
}
