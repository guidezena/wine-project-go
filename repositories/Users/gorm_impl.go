package users

import (
	"context"
	"db-go/login/entities"

	"gorm.io/gorm"
)

type userDBGorm struct {
	writer *gorm.DB
	reader *gorm.DB
}

func NewGorm(w *gorm.DB, r *gorm.DB) UserRepositoryInterface {
	return &userDBGorm{writer: w, reader: r}
}

func (u *userDBGorm) GetAll(ctx context.Context) ([]entities.User, error) {
	var users []entities.User

	u.reader.Table("users").Find(&users)

	if u.reader.Error != nil {
		return nil, u.reader.Error
	}

	return users, nil
}
