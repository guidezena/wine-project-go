package entities

type User struct {
	IdUser int    `db:"iduser" json:"iduser"`
	Name   string `db:"name" json:"name"`
	Email  string `db:"email" json:"email"`
}
