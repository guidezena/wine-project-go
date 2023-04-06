package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string   `gorm:"name" json:"name"`
	Email    string   `gorm:"email" json:"email"`
	Password string   `gorm:"password" json:"password"`
	UserType UserType `gorm:"user_type" json:"user_type"`
}
