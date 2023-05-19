package dishes

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

func AddDishHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request AddDish")

	var dish entities.Dish
	err := json.NewDecoder(r.Body).Decode(&dish)

	log.Printf(dish.Name)

	if err != nil {
		utils.SendError(w, "Decode error", http.StatusBadRequest)
		return
	}

	writer := dbConnection.GetWriterGorm()
	errorToWrite := createDish(writer, dish)
	dbConnection.CloseDbConnection(writer)

	if errorToWrite != nil {
		log.Printf("errorToWrite")

		utils.SendError(w, errorToWrite.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]string{
		"message": "Prato criada com sucesso!",
	}

	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
	log.Printf("Prato criada com sucesso!")
}

func createDish(db *gorm.DB, dish entities.Dish) error {
	log.Printf("createDish")

	newDish := entities.Dish{
		CategoryId:   dish.CategoryId,
		RestaurantId: dish.RestaurantId,
		Name:         dish.Name,
		Image:        dish.Image,
		Description:  dish.Description,
	}

	result := db.Create(&newDish)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &entities.CustomError{Message: "Nenhuma linha foi afetada"}
	}

	return nil
}

func GetDishesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request GetDishes")

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Erro ao analisar os parâmetros da solicitação", http.StatusBadRequest)
		return
	}
	reader := dbConnection.GetReaderGorm()

	idRestaurant := r.Form.Get("restaurant_id")
	idCategoria := r.Form.Get("category_id")

	var dishes []entities.Dish
	query := reader

	if idRestaurant != "" {
		log.Printf("Buscando restaurant_id")
		query = query.Where("restaurant_id = ?", idRestaurant)
	}
	if idCategoria != "" {
		log.Printf("Buscando category_id")
		query = query.Where("category_id = ?", idCategoria)
	}

	err = query.Find(&dishes).Error

	dbConnection.CloseDbConnection(query)

	if err != nil {
		http.Error(w, "Erro ao obter pratos", http.StatusInternalServerError)
		return
	}

	jsonDishes, err := json.Marshal(dishes)
	if err != nil {
		http.Error(w, "Erro ao converter pratos para JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonDishes)
}

func GetDishHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request GetDishHandler")

	dishID := mux.Vars(r)["id"]

	reader := dbConnection.GetReaderGorm()
	errorToDelete := getDish(reader, dishID)
	dbConnection.CloseDbConnection(reader)

	if errorToDelete != nil {
		utils.SendError(w, errorToDelete.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"message": "sucesso!",
	}

	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	log.Printf("sucesso!")
}

func getDish(db *gorm.DB, dishID string) error {
	log.Printf("getDish")

	var dish entities.Dish
	result := db.First(&dish, dishID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &entities.CustomError{Message: "Nenhuma linha foi afetada"}
	}

	return nil
}

func DeleteDishHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request DeleteDishHandler")

	dishID := mux.Vars(r)["id"]

	writer := dbConnection.GetWriterGorm()
	errorToDelete := deleteDish(writer, dishID)
	dbConnection.CloseDbConnection(writer)

	if errorToDelete != nil {
		utils.SendError(w, errorToDelete.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"message": "Prato excluído com sucesso!",
	}

	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	log.Printf("Prato excluído com sucesso!")
}

func deleteDish(db *gorm.DB, dishID string) error {
	log.Printf("deleteDish")

	result := db.Delete(&entities.Dish{}, dishID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &entities.CustomError{Message: "Nenhuma linha foi afetada"}
	}

	return nil
}

func UpdateDishHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request UpdateDishHandler")

	dishID := mux.Vars(r)["id"]

	var updatedDish entities.Dish
	err := json.NewDecoder(r.Body).Decode(&updatedDish)

	if err != nil {
		utils.SendError(w, "Decode error", http.StatusBadRequest)
		return
	}

	writer := dbConnection.GetWriterGorm()
	errorToUpdate := updateDish(writer, dishID, updatedDish)
	dbConnection.CloseDbConnection(writer)

	if errorToUpdate != nil {
		utils.SendError(w, errorToUpdate.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"message": "Prato atualizado com sucesso!",
	}

	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	log.Printf("Prato atualizado com sucesso!")
}

func updateDish(db *gorm.DB, dishID string, updatedDish entities.Dish) error {
	log.Printf("updateDish")

	result := db.Model(&entities.Dish{}).Where("ID = ?", dishID).Updates(updatedDish)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &entities.CustomError{Message: "Nenhuma linha foi afetada"}
	}

	return nil
}
