package entities

import "gorm.io/gorm"

type UserForm struct {
	gorm.Model
	Name            string `gorm:"name" json:"name"`
	Email           string `gorm:"email" json:"email"`
	Password        string `gorm:"password" json:"password"`
	ConfirmPassword string `gorm:"password_confirmation" json:"password_confirmation"`
	IsAdmin         bool   `gorm:"is_admin" json:"is_admin"`
}

type User struct {
	gorm.Model
	Name     string `gorm:"name" json:"name"`
	Email    string `gorm:"email" json:"email"`
	Password string `gorm:"password" json:"password"`
	IsAdmin  bool   `gorm:"is_admin" json:"is_admin"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
