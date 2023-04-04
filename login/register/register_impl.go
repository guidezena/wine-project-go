package register

import (
	"db-go/configs"
	"db-go/login/entities"
	"fmt"
	"log"
	"net/http"

	"gorm.io/gorm"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Recebendo requisição em CreateUserHandler")

	if r.Method != "POST" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Ler o corpo da requisição e decodificar os dados do novo usuário
	//var user entities.User
	user, err := parseUser(r)

	if err != nil {
		http.Error(w, "decode error", http.StatusBadRequest)
		return
	}

	// Verificar se o usuário já existe (por exemplo, pelo email)

	reader := configs.GetReaderGorm()
	existsEmail, err := UserExistsByEmail(reader, user.Email)

	if err != nil {
		http.Error(w, "OUTRO ERROR", http.StatusBadRequest)
		return
	}

	if existsEmail {
		http.Error(w, "Este email já está em uso", http.StatusConflict)
		return
	}

	// Se não existir, criar um novo usuário no banco de dados

	writer := configs.GetWriterGorm()
	errorToWrite := addUser(writer, *user)

	if errorToWrite != nil {
		http.Error(w, errorToWrite.Error(), http.StatusBadRequest)
		return
	}

	// Retornar o novo usuário criado com um status HTTP 201 (criado)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Usuário criado com sucesso!")
}

func UserExistsByEmail(db *gorm.DB, email string) (bool, error) {
	var count int64
	err := db.Unscoped().Model(&entities.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func addUser(db *gorm.DB, user entities.User) error {
	result := db.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &entities.CustomError{Message: "Nenhuma linha foi afetada"}
	}

	return nil
}

func parseUser(r *http.Request) (*entities.User, error) {
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	userType := r.FormValue("userType")

	// criando um objeto do tipo UserType
	var ut entities.UserType
	if userType == "admin" {
		ut = entities.Admin
	} else if userType == "basic" {
		ut = entities.Basic
	} else {
		return nil, fmt.Errorf("invalid userType")
	}

	// criando um objeto do tipo User com as informações da request
	user := &entities.User{
		Name:     name,
		Email:    email,
		Password: password,
		UserType: ut,
	}

	return user, nil
}
