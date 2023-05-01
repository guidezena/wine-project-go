package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"wine-project-go/login/entities"
	"wine-project-go/login/register"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func GenerateToken(user *entities.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"nbf":      time.Now().Unix(),
		"iat":      time.Now().Unix(),
		"sub":      user.Email,
		"username": user.Name,
	})
	return token.SignedString([]byte("my-secret-key"))
}

func Authenticate(email string, password string) (*entities.User, error) {
	user, err := register.GetUser(email)

	// Comparar a senha hash armazenada com a senha fornecida pelo usuário
	log.Printf(user.Email)
	log.Printf(user.Password)

	log.Printf(email)
	log.Printf(password)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, fmt.Errorf("Invalid credentials")
		} else {
			return nil, err
		}
	}

	return user, nil
}

func ValidateToken(tokenString string) error {
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
