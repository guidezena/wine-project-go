package entities

import "gorm.io/gorm"

type Drink struct {
	gorm.Model
	Name        string `gorm:"name" json:"name"`
	Image       string `gorm:"image" json:"image"`
	Description string `gorm:"description" json:"description"`
}

func (Drink) TableName() string {
	return "drinks"
}
