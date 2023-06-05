package categories

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

func AddCategoryHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request AddCategory")

	var category entities.Category
	err := json.NewDecoder(r.Body).Decode(&category)

	log.Printf(category.Name)

	if err != nil {
		utils.SendError(w, "Decode error", http.StatusBadRequest)
		return
	}

	writer := dbConnection.GetWriterGorm()
	errorToWrite := createCategory(writer, category)
	dbConnection.CloseDbConnection(writer)

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

func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request GetCategories")

	reader := dbConnection.GetReaderGorm()
	categories, err := getCategories(reader)
	dbConnection.CloseDbConnection(reader)

	if err != nil {
		http.Error(w, "Erro ao obter categorias", http.StatusInternalServerError)
		return
	}

	jsonCategories, err := json.Marshal(categories)
	if err != nil {
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

func GetCategoriesForIdHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request GetCategoriesForIdHandler")

	categoryID := mux.Vars(r)["id"]

	reader := dbConnection.GetReaderGorm()
	categories, err := getCategoriesForId(reader, categoryID)
	dbConnection.CloseDbConnection(reader)

	if err != nil {
		http.Error(w, "Erro ao obter categorias", http.StatusInternalServerError)
		return
	}

	jsonCategories, err := json.Marshal(categories)
	if err != nil {
		http.Error(w, "Erro ao converter categorias para JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonCategories)
}

func getCategoriesForId(db *gorm.DB, categoryId string) (*entities.Category, error) {

	var category entities.Category
	result := db.First(&category, categoryId)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, &entities.CustomError{Message: "Nenhuma linha foi afetada"}
	}

	return &category, nil

}

func DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request DeleteCategoryHandler")

	categoryID := mux.Vars(r)["id"]

	writer := dbConnection.GetWriterGorm()
	errorToDelete := deleteCategory(writer, categoryID)
	dbConnection.CloseDbConnection(writer)

	if errorToDelete != nil {
		utils.SendError(w, errorToDelete.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"message": "Categoria excluída com sucesso!",
	}

	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	log.Printf("Categoria excluída com sucesso!")
}

func deleteCategory(db *gorm.DB, categoryID string) error {
	log.Printf("deleteCategory")

	result := db.Delete(&entities.Category{}, categoryID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &entities.CustomError{Message: "Nenhuma linha foi afetada"}
	}

	return nil
}

func UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request UpdateCategoryHandler")

	categoryID := mux.Vars(r)["id"]

	var updatedCategory entities.Category
	err := json.NewDecoder(r.Body).Decode(&updatedCategory)

	if err != nil {
		utils.SendError(w, "Decode error", http.StatusBadRequest)
		return
	}

	writer := dbConnection.GetWriterGorm()
	errorToUpdate := updateCategory(writer, categoryID, updatedCategory)
	dbConnection.CloseDbConnection(writer)

	if errorToUpdate != nil {
		utils.SendError(w, errorToUpdate.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"message": "Categoria atualizada com sucesso!",
	}

	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	log.Printf("Categoria atualizada com sucesso!")
}

func updateCategory(db *gorm.DB, categoryID string, updatedCategory entities.Category) error {
	log.Printf("updateCategory")

	result := db.Model(&entities.Category{}).Where("ID = ?", categoryID).Updates(updatedCategory)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &entities.CustomError{Message: "Nenhuma linha foi afetada"}
	}

	return nil
}
