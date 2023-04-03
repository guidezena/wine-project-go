package main

import (
	"context"
	"db-go/configs"
	"db-go/repositories"
	"fmt"
	"log"
)

func main() {
	repo := repositories.New(repositories.Options{
		WriterSqlx: configs.GetWriterSqlx(),
		ReaderSqlx: configs.GetReaderSqlx(),
		WriterGorm: configs.GetWriterGorm(),
		ReaderGorm: configs.GetReaderGorm(),
	})

	users, err := repo.User.GetAll(context.Background())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(users)
}
