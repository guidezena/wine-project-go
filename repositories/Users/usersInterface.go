package users

import (
	"context"
	"db-go/entities"
)

type UserRepositoryInterface interface {
	GetAll(ctx context.Context) ([]entities.User, error)
}
