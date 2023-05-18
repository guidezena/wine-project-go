package users

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

func GetUsers(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request GetUsers")

	reader := dbConnection.GetReaderGorm()
	drinks, err := getUsers(reader)
	dbConnection.CloseDbConnection(reader)

	if err != nil {
		http.Error(w, "Erro ao obter users", http.StatusInternalServerError)
		return
	}

	jsonUsers, err := json.Marshal(drinks)
	if err != nil {
		http.Error(w, "Erro ao converter users para JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonUsers)
}

func getUsers(db *gorm.DB) ([]entities.User, error) {
	var users []entities.User
	err := db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request DeleteUser")

	userID := mux.Vars(r)["id"]

	writer := dbConnection.GetWriterGorm()
	errorToDelete := deleteUser(writer, userID)
	dbConnection.CloseDbConnection(writer)

	if errorToDelete != nil {
		utils.SendError(w, errorToDelete.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"message": "Usuario excluído com sucesso!",
	}

	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	log.Printf("Usuario excluída com sucesso!")
}

func deleteUser(db *gorm.DB, userID string) error {
	log.Printf("deleteUser")

	result := db.Delete(&entities.User{}, userID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &entities.CustomError{Message: "Nenhuma linha foi afetada"}
	}

	return nil
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request UpdateUser")

	userID := mux.Vars(r)["id"]

	var updatedUser entities.User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)

	if err != nil {
		utils.SendError(w, "Decode error", http.StatusBadRequest)
		return
	}

	writer := dbConnection.GetWriterGorm()
	errorToUpdate := updateUser(writer, userID, updatedUser)
	dbConnection.CloseDbConnection(writer)

	if errorToUpdate != nil {
		utils.SendError(w, errorToUpdate.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"message": "Usuario atualizado com sucesso!",
	}

	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	log.Printf("Usuario atualizado com sucesso!")
}

func updateUser(db *gorm.DB, userID string, updatedUser entities.User) error {
	log.Printf("updateUser")

	result := db.Model(&entities.User{}).Where("ID = ?", userID).Updates(updatedUser)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &entities.CustomError{Message: "Nenhuma linha foi afetada"}
	}

	return nil
}
