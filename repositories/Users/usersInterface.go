package users

import (
	"context"
	"wine-project-go/login/entities"
)

type UserRepositoryInterface interface {
	GetAll(ctx context.Context) ([]entities.User, error)
}
