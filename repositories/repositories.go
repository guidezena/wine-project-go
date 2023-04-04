package repositories

import (
	users "db-go/repositories/Users"

	"gorm.io/gorm"
)

type Options struct {
	WriterGorm *gorm.DB
	ReaderGorm *gorm.DB
}

type Container struct {
	User users.UserRepositoryInterface
}

func New(options Options) *Container {
	return &Container{
		User: users.NewGorm(options.WriterGorm, options.ReaderGorm),
	}
}
