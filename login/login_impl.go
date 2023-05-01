package login

import (
	"encoding/json"
	"log"
	"net/http"
	"wine-project-go/login/entities"
	"wine-project-go/login/auth"
)

func Login(w http.ResponseWriter, r *http.Request) {

	log.Printf("Receiving request in Login")

	credentials := parseCredentials(r)

	user, err := session.Authenticate(credentials.Email, credentials.Password)
	if err != nil {
		log.Printf("autentication")

		http.Error(w, err.Error(), http.StatusUnauthorized)

		return
	}

	tokenString, err := session.GenerateToken(user)
	if err != nil {
		log.Printf("generate token")
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func parseCredentials(r *http.Request) *entities.Credentials {
	email := r.FormValue("email")
	password := r.FormValue("password")

	credentials := &entities.Credentials{
		Email:    email,
		Password: password,
	}

	return credentials
}
