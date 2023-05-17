package entities

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name  string `gorm:"name" json:"name"`
	Image string `gorm:"image" json:"image"`
}

func (Category) TableName() string {
	return "categories"
}
