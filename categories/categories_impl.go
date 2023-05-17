package categories

import (
	"encoding/json"
	"log"
	"net/http"
	"wine-project-go/configs"
	"wine-project-go/entities"
	"wine-project-go/utils"

	"gorm.io/gorm"
)

func AddCategory(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request AddCategory")

	var category entities.Category
	err := json.NewDecoder(r.Body).Decode(&category)

	log.Printf(category.Name)

	if err != nil {
		utils.SendError(w, "Decode error", http.StatusBadRequest)
		return
	}

	writer := configs.GetWriterGorm()
	errorToWrite := createCategory(writer, category)
	configs.CloseDbConnection(writer)

	if errorToWrite != nil {
		log.Printf("errorToWrite")

		utils.SendError(w, errorToWrite.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]string{
		"message": "Categoria criada com sucesso!",
	}

	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
	log.Printf("Categoria criada com sucesso!")
}

func createCategory(db *gorm.DB, category entities.Category) error {
	log.Printf("createCategory")

	newCategory := entities.Category{
		Name:  category.Name,
		Image: category.Image,
	}

	result := db.Create(&newCategory)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &entities.CustomError{Message: "Nenhuma linha foi afetada"}
	}

	return nil
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request GetCategories")

	reader := configs.GetReaderGorm()
	categories, err := getCategories(reader)
	configs.CloseDbConnection(reader)

	if err != nil {
		// Trate o erro
		http.Error(w, "Erro ao obter categorias", http.StatusInternalServerError)
		return
	}

	// Enviar as categorias como resposta JSON
	jsonCategories, err := json.Marshal(categories)
	if err != nil {
		// Trate o erro
		http.Error(w, "Erro ao converter categorias para JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonCategories)
}

func getCategories(db *gorm.DB) ([]entities.Category, error) {
	var categories []entities.Category
	err := db.Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}
