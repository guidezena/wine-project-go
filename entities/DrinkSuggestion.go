package entities

import "gorm.io/gorm"

type DrinkSuggestion struct {
	gorm.Model
	DishID  int   `gorm:"dish_id" json:"dish_id"`
	DrinkID int   `gorm:"drink_id" json:"drink_id"`
	Drink   Drink `gorm:"foreignKey:DrinkID"`
}

func (DrinkSuggestion) TableName() string {
	return "drink_suggestions"
}
