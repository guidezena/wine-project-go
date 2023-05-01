package register

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"wine-project-go/configs"
	"wine-project-go/login/entities"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request in CreateUserHandler")

	if r.Method != "POST" {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userForm entities.UserForm
	err := json.NewDecoder(r.Body).Decode(&userForm)
	// Ler o corpo da requisição e decodificar os dados do novo usuário
	//var user entities.User
	//user, err := parseUser(r)
	log.Printf(userForm.Password)
	log.Printf(userForm.ConfirmPassword)

	if userForm.Password != userForm.ConfirmPassword {
		sendError(w, "As senhas nao sao iguais", http.StatusBadRequest)
		return
	}

	if err != nil {
		sendError(w, "Decode error", http.StatusBadRequest)
		return
	}

	// Verificar se o usuário já existe (por exemplo, pelo email)

	reader := configs.GetReaderGorm()
	existsEmail, err := UserExistsByEmail(reader, userForm.Email)

	if err != nil {
		sendError(w, "Error to validate if the registered email already exists", http.StatusBadRequest)
		return
	}

	if existsEmail {
		sendError(w, "This email is already registered", http.StatusConflict)
		return
	}

	log.Printf("print 1")

	// Se não existir, criar um novo usuário no banco de dados

	writer := configs.GetWriterGorm()
	errorToWrite := addUser(writer, userForm)

	if errorToWrite != nil {
		log.Printf("errorToWrite")

		sendError(w, errorToWrite.Error(), http.StatusBadRequest)
		return
	}

	// Retornar o novo usuário criado com um status HTTP 201 (criado)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Usuário criado com sucesso!")
	log.Printf("Usuário criado com sucesso!")

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

func addUser(db *gorm.DB, userForm entities.UserForm) error {
	log.Printf("Add user")

	user := entities.User{
		Name:     userForm.Name,
		Email:    userForm.Email,
		Password: userForm.Password,
	}

	result := db.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &entities.CustomError{Message: "Nenhuma linha foi afetada"}
	}

	return nil
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	log.Printf(message)

	w.WriteHeader(statusCode)
	//w.Header().Set("X-Status-Message", message)
	fmt.Fprintf(w, message)
}

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}

func GetUser(email string) (*entities.User, error) {
	log.Printf("getUser")

	var user entities.User

	reader := configs.GetReaderGorm()

	log.Printf(email)

	// Fazer uma consulta para buscar o usuário pelo email
	err := reader.Where(&entities.User{Email: email}).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func parseUser(r *http.Request) (*entities.User, error) {
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("password_confirmation")

	if password != passwordConfirm {
		return nil, fmt.Errorf("Your passwords do not match")
	}

	log.Printf(hashPassword(password))

	// criando um objeto do tipo User com as informações da request
	user := &entities.User{
		Name:     name,
		Email:    email,
		Password: hashPassword(password),
	}

	return user, nil
}
