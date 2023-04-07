package session

import (
	"db-go/login/entities"
	"db-go/login/register"
	"fmt"
	"log"
	"time"

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

func Authenticate(email, password string) (*entities.User, error) {
	user, err := register.GetUser(email)

	// Comparar a senha hash armazenada com a senha fornecida pelo usuário
	log.Printf(user.Password)
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
