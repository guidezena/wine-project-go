package suggestions

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

func AddDrinkSuggestionHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request AddDrinkSuggestion")

	var drinkSuggestion entities.DrinkSuggestion
	err := json.NewDecoder(r.Body).Decode(&drinkSuggestion)

	if err != nil {
		utils.SendError(w, "Decode error", http.StatusBadRequest)
		return
	}

	writer := dbConnection.GetWriterGorm()
	errorToWrite := createDrinkSuggestion(writer, drinkSuggestion)
	dbConnection.CloseDbConnection(writer)

	if errorToWrite != nil {
		log.Printf("errorToWrite")

		utils.SendError(w, errorToWrite.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]string{
		"message": "Sugestão criada com sucesso!",
	}

	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
	log.Printf("Sugestão criada com sucesso!")
}

func createDrinkSuggestion(db *gorm.DB, drinkSuggestion entities.DrinkSuggestion) error {
	log.Printf("createDrinkSuggestion")

	newDrinkSuggestion := entities.DrinkSuggestion{
		DishID:  drinkSuggestion.DishID,
		DrinkID: drinkSuggestion.DrinkID,
	}

	result := db.Create(&newDrinkSuggestion)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &entities.CustomError{Message: "Nenhuma linha foi afetada"}
	}

	return nil
}

func GetDrinkSuggestionsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request GetDrinkSuggestions")

	dishID := mux.Vars(r)["dishID"]

	// Verifique se o usuário é premium
	isPremium := true // Defina essa variável de acordo com a lógica do seu aplicativo

	reader := dbConnection.GetReaderGorm()
	drinkSuggestions, err := getDrinkSuggestions(reader, dishID, isPremium)
	dbConnection.CloseDbConnection(reader)

	if err != nil {
		utils.SendError(w, "Erro ao obter drink suggestions", http.StatusInternalServerError)
		return
	}

	jsonDrinkSuggestions, err := json.Marshal(drinkSuggestions)
	if err != nil {
		http.Error(w, "Erro ao converter drink suggestions para JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonDrinkSuggestions)
}

func getDrinkSuggestions(db *gorm.DB, dishID string, isPremium bool) ([]entities.DrinkSuggestion, error) {
	var drinkSuggestions []entities.DrinkSuggestion

	query := db.Where("dish_id = ?", dishID)

	if isPremium {
		query = query.Limit(3)
	} else {
		query = query.Limit(1)
	}

	err := query.Preload("Drink").Find(&drinkSuggestions).Error
	if err != nil {
		return nil, err
	}

	return drinkSuggestions, nil
}

func DeleteDrinkSuggestionHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request DeleteDrinkSuggestionHandler")

	params := mux.Vars(r)
	dishID := params["dishID"]
	drinkID := params["drinkID"]

	writer := dbConnection.GetWriterGorm()
	errorToDelete := deleteDrinkSuggestion(writer, dishID, drinkID)
	dbConnection.CloseDbConnection(writer)

	if errorToDelete != nil {
		utils.SendError(w, errorToDelete.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"message": "Sugestão excluída com sucesso!",
	}

	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	log.Printf("Sugestão excluída com sucesso!")
}

func deleteDrinkSuggestion(db *gorm.DB, dishID string, drinkID string) error {
	log.Printf("deleteDrinkSuggestion")

	var suggestion entities.DrinkSuggestion
	result := db.Where("dish_id = ? AND drink_id = ?", dishID, drinkID).Delete(&suggestion)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &entities.CustomError{Message: "Nenhuma linha foi afetada"}
	}

	return nil
}
