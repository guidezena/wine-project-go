package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"wine-project-go/entities"
	"wine-project-go/repositories/register"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func GenerateToken(user *entities.User) (string, error) {
	log.Printf("GenerateToken")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"nbf":      time.Now().Unix(),
		"iat":      time.Now().Unix(),
		"sub":      user.Email,
		"username": user.Name,
		"is_admin": user.IsAdmin,
	})
	return token.SignedString([]byte("my-secret-key"))
}

func Authenticate(email string, password string) (*entities.User, error) {
	log.Printf("Authenticate")

	user, err := register.GetUser(email)

	// Comparar a senha hash armazenada com a senha fornecida pelo usuário
	log.Printf(user.Email)
	log.Printf(user.Password)

	log.Printf(email)
	log.Printf(password)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			log.Printf("Invalid credentials")
			return nil, fmt.Errorf("Invalid credentials")
		} else {
			log.Printf("Erro CompareHashAndPassword")
			return nil, err
		}
	}

	return user, nil
}

func ValidateToken(tokenString string) error {
	log.Printf("ValidateToken")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("my-secret-key"), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("Invalid token")
	}

	return nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		err := ValidateToken(strings.Replace(tokenString, "Bearer ", "", 1))
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
