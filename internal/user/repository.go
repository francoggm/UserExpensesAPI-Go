package user

import (
	"database/sql"
)

type repository struct {
	db *sql.DB
}

// returns repository struct implementing the Repository interface
func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateUser(user *User) (int64, error) {
	var id int64

	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id"

	err := r.db.QueryRow(query, user.Email, user.Password).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repository) GetUserByEmail(email string) (user *User, err error) {
	return
}
