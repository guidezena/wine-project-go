package users

import (
	"context"
	"db-go/entities"

	"github.com/jmoiron/sqlx"
)

type userDBSqlx struct {
	writer *sqlx.DB
	reader *sqlx.DB
}

func NewSqlx(w *sqlx.DB, r *sqlx.DB) UserRepositoryInterface {
	return &userDBSqlx{writer: w, reader: r}
}

func (u *userDBSqlx) GetAll(ctx context.Context) ([]entities.User, error) {
	var users []entities.User
	err := u.reader.SelectContext(ctx, &users, `
		select * from users
	`)

	if err != nil {
		return nil, err
	}

	return users, nil
}
