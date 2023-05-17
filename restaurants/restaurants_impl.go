package restaurants

import (
	"encoding/json"
	"log"
	"net/http"
	"wine-project-go/configs"
	"wine-project-go/entities"
	"wine-project-go/utils"

	"gorm.io/gorm"
)

func AddRestaurant(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request AddRestaurant")

	var restaurant entities.Restaurant
	err := json.NewDecoder(r.Body).Decode(&restaurant)

	log.Printf(restaurant.Name)

	if err != nil {
		utils.SendError(w, "Decode error", http.StatusBadRequest)
		return
	}

	writer := configs.GetWriterGorm()
	errorToWrite := createRestaurant(writer, restaurant)
	configs.CloseDbConnection(writer)

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

func GetRestaurants(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request GetRestaurants")

	reader := configs.GetReaderGorm()
	restaurants, err := getRestaurants(reader)
	configs.CloseDbConnection(reader)

	if err != nil {
		// Trate o erro
		http.Error(w, "Erro ao obter restaurantes", http.StatusInternalServerError)
		return
	}

	// Enviar os restaurantes como resposta JSON
	jsonRestaurants, err := json.Marshal(restaurants)
	if err != nil {
		// Trate o erro
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
