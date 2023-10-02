package user

import (
	"database/sql"
	"time"
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

func (r *repository) CreateUser(user *User) error {
	var id int64

	query := "INSERT INTO users (email, password, name, created_at) VALUES ($1, $2, $3, $4) RETURNING id"

	createdDate := time.Now().Format("2006-01-02 15:04:05")

	err := r.db.QueryRow(query, user.Email, user.Password, user.Name, createdDate).Scan(&id)
	if err != nil {
		return err
	}

	user.ID = id
	user.CreatedAt = createdDate

	return nil
}

func (r *repository) GetUserByEmail(email string) (*User, error) {
	var user User
	var lastLogin sql.NullString

	err := r.db.QueryRow("SELECT * FROM users WHERE email=$1", email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &lastLogin)
	if err != nil {
		return nil, err
	}

	user.LastLogin = lastLogin.String

	return &user, nil
}
