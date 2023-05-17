package login

import (
	"encoding/json"
	"log"
	"net/http"
	"wine-project-go/entities"
	"wine-project-go/login/auth"
)

func Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving request in Login")

	var creds entities.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf(creds.Email)
	log.Printf(creds.Password)

	user, err := auth.Authenticate(creds.Email, creds.Password)

	if err != nil {
		http.Error(w, "StatusUnauthorized", http.StatusUnauthorized)
		return
	}

	tokenString, err := auth.GenerateToken(user)
	if err != nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)

		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
