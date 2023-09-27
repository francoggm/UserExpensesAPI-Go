package user

import (
	"database/sql"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateUser(email string, password string) (id int64, err error) {
	return
}

func (r *repository) GetUserByEmail(email string) (user *User, err error) {
	return
}
