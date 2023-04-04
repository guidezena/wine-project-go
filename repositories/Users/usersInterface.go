package users

import (
	"context"
	"db-go/login/entities"
)

type UserRepositoryInterface interface {
	GetAll(ctx context.Context) ([]entities.User, error)
}
