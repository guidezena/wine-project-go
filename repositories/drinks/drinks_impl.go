package drinks

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

func AddDrinkHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request AddDrink")

	var drink entities.Drink
	err := json.NewDecoder(r.Body).Decode(&drink)

	if err != nil {
		utils.SendError(w, "Decode error", http.StatusBadRequest)
		return
	}

	writer := dbConnection.GetWriterGorm()
	errorToWrite := createDrink(writer, drink)
	dbConnection.CloseDbConnection(writer)

	if errorToWrite != nil {
		log.Printf("errorToWrite")

		utils.SendError(w, errorToWrite.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]string{
		"message": "Drink criado com sucesso!",
	}

	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
	log.Printf("Drink criada com sucesso!")
}

func createDrink(db *gorm.DB, drink entities.Drink) error {
	log.Printf("createDrink")

	newDrink := entities.Drink{
		Name:        drink.Name,
		Image:       drink.Image,
		Description: drink.Description,
	}

	result := db.Create(&newDrink)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &entities.CustomError{Message: "Nenhuma linha foi afetada"}
	}

	return nil
}

func GetDrinksHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request GetDrinks")

	reader := dbConnection.GetReaderGorm()
	drinks, err := getDrinks(reader)
	dbConnection.CloseDbConnection(reader)

	if err != nil {
		http.Error(w, "Erro ao obter drinks", http.StatusInternalServerError)
		return
	}

	jsonDrinks, err := json.Marshal(drinks)
	if err != nil {
		http.Error(w, "Erro ao converter drinks para JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonDrinks)
}

func getDrinks(db *gorm.DB) ([]entities.Drink, error) {
	var drinks []entities.Drink
	err := db.Find(&drinks).Error
	if err != nil {
		return nil, err
	}
	return drinks, nil
}

func DeleteDrinkHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request DeleteDrinkHandler")

	drinkID := mux.Vars(r)["id"]

	writer := dbConnection.GetWriterGorm()
	errorToDelete := deleteDrink(writer, drinkID)
	dbConnection.CloseDbConnection(writer)

	if errorToDelete != nil {
		utils.SendError(w, errorToDelete.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"message": "Drink excluído com sucesso!",
	}

	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	log.Printf("Drink excluído com sucesso!")
}

func deleteDrink(db *gorm.DB, drinkID string) error {
	log.Printf("deleteDrink")

	result := db.Delete(&entities.Drink{}, drinkID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &entities.CustomError{Message: "Nenhuma linha foi afetada"}
	}

	return nil
}

func UpdateDrinkHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request UpdateDrinkHandler")

	drinkID := mux.Vars(r)["id"]

	var updatedDrink entities.Drink
	err := json.NewDecoder(r.Body).Decode(&updatedDrink)

	if err != nil {
		utils.SendError(w, "Decode error", http.StatusBadRequest)
		return
	}

	writer := dbConnection.GetWriterGorm()
	errorToUpdate := updateDrink(writer, drinkID, updatedDrink)
	dbConnection.CloseDbConnection(writer)

	if errorToUpdate != nil {
		utils.SendError(w, errorToUpdate.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"message": "Drink atualizado com sucesso!",
	}

	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	log.Printf("Drink atualizado com sucesso!")
}

func updateDrink(db *gorm.DB, drinkID string, updatedDrink entities.Drink) error {
	log.Printf("updateDrink")

	result := db.Model(&entities.Drink{}).Where("ID = ?", drinkID).Updates(updatedDrink)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &entities.CustomError{Message: "Nenhuma linha foi afetada"}
	}

	return nil
}
