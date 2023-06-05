package auth

import (
	"fmt"
	"log"
	"time"
	"wine-project-go/entities"
	"wine-project-go/repositories/register"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func GenerateToken(user *entities.User) (string, error) {
	log.Printf("GenerateToken")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    user.ID,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
		"nbf":        time.Now().Unix(),
		"iat":        time.Now().Unix(),
		"sub":        user.Email,
		"username":   user.Name,
		"is_admin":   user.IsAdmin,
		"is_premium": user.IsPremium,
	})
	return token.SignedString([]byte("my-secret-key"))
}

func Authenticate(email string, password string) (*entities.User, error) {
	log.Printf("Authenticate")

	user, err := register.GetUser(email)

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
