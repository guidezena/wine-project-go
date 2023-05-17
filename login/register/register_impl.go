package register

import (
	"encoding/json"
	"log"
	"net/http"
	"wine-project-go/configs"
	"wine-project-go/login/entities"
	"wine-project-go/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request in CreateUserHandler")

	if r.Method != "POST" {
		utils.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
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
		utils.SendError(w, "As senhas nao sao iguais", 401)
		return
	}

	if err != nil {
		utils.SendError(w, "Decode error", http.StatusBadRequest)
		return
	}

	// Verificar se o usuário já existe (por exemplo, pelo email)

	reader := configs.GetReaderGorm()
	existsEmail, err := UserExistsByEmail(reader, userForm.Email)
	configs.CloseDbConnection(reader)

	if err != nil {
		utils.SendError(w, "Error to validate if the registered email already exists", http.StatusBadRequest)
		return
	}

	if existsEmail {
		utils.SendError(w, "This email is already registered", http.StatusConflict)
		return
	}

	// Se não existir, criar um novo usuário no banco de dados

	writer := configs.GetWriterGorm()
	errorToWrite := addUser(writer, userForm)
	configs.CloseDbConnection(writer)

	if errorToWrite != nil {
		log.Printf("errorToWrite")

		utils.SendError(w, errorToWrite.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]string{
		"message": "Usuário criado com sucesso!",
	}

	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
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
		Password: hashPassword(userForm.Password),
		IsAdmin:  userForm.IsAdmin,
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
	err := reader.Where(&entities.User{Email: email}).First(&user).Error
	configs.CloseDbConnection(reader)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
