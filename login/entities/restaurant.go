package entities

import "gorm.io/gorm"

type Restaurant struct {
	gorm.Model
	UserId      string `gorm:"id_user" json:"id_user"`
	Name        string `gorm:"name" json:"name"`
	Image       string `gorm:"image" json:"image"`
	Description string `gorm:"description" json:"description"`
	Address     string `gorm:"address" json:"address"`
}

func (Restaurant) TableName() string {
	return "restaurants"
}
