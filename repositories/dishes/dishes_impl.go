package dishes

import (
	"encoding/json"
	"log"
	"net/http"
	"wine-project-go/dbConnection"
	"wine-project-go/entities"
	"wine-project-go/utils"

	"gorm.io/gorm"
)

func AddDish(w http.ResponseWriter, r *http.Request) {
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

func GetDishes(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request GetDishes")

	// Obtenha os parâmetros da solicitação POST
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Erro ao analisar os parâmetros da solicitação", http.StatusBadRequest)
		return
	}
	reader := dbConnection.GetReaderGorm()

	// Verifique os parâmetros fornecidos
	idRestaurant := r.Form.Get("restaurant_id")
	idCategoria := r.Form.Get("category_id")

	// Realize a busca com base nos parâmetros fornecidos
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

	log.Printf("DEU BOM A BUSCA")

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonDishes)
}