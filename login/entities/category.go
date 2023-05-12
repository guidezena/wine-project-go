package entities

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name string `gorm:"name" json:"name"`
}

func (Category) TableName() string {
	return "categories"
}
