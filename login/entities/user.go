package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name            string `gorm:"name" json:"name"`
	Email           string `gorm:"email" json:"email"`
	Password        string `gorm:"password" json:"password"`
	ConfirmPassword string `gorm:"password_confirmation" json:"password_confirmation"`
}

type DBUser struct {
	gorm.Model
	Name     string `gorm:"name" json:"name"`
	Email    string `gorm:"email" json:"email"`
	Password string `gorm:"password" json:"password"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
