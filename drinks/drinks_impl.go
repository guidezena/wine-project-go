package drinks

import (
	"encoding/json"
	"log"
	"net/http"
	"wine-project-go/configs"
	"wine-project-go/login/entities"
	"wine-project-go/utils"

	"gorm.io/gorm"
)

func AddDrink(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request AddDrink")

	var drink entities.Drink
	err := json.NewDecoder(r.Body).Decode(&drink)

	log.Printf(drink.Name)

	if err != nil {
		utils.SendError(w, "Decode error", http.StatusBadRequest)
		return
	}

	writer := configs.GetWriterGorm()
	errorToWrite := createDrink(writer, drink)

	if errorToWrite != nil {
		log.Printf("errorToWrite")

		utils.SendError(w, errorToWrite.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]string{
		"message": "Drink criada com sucesso!",
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

func GetDrinks(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request GetDrinks")

	reader := configs.GetReaderGorm()
	drinks, err := getDrinks(reader)

	if err != nil {
		// Trate o erro
		http.Error(w, "Erro ao obter drinks", http.StatusInternalServerError)
		return
	}

	// Enviar as categorias como resposta JSON
	jsonDrinks, err := json.Marshal(drinks)
	if err != nil {
		// Trate o erro
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
