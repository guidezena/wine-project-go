package entities

import "gorm.io/gorm"

type Dish struct {
	gorm.Model
	CategoryId   int    `gorm:"category_id" json:"category_id"`
	RestaurantId int    `gorm:"restaurant_id" json:"restaurant_id"`
	Name         string `gorm:"name" json:"name"`
	Image        string `gorm:"image" json:"image"`
	Description  string `gorm:"description" json:"description"`
}

func (Dish) TableName() string {
	return "dishes"
}
