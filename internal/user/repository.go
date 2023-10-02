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

func (r *repository) CreateUser(user *User) (int64, error) {
	var id int64

	query := "INSERT INTO users (email, password, created_at) VALUES ($1, $2, $3) RETURNING user_id"

	err := r.db.QueryRow(query, user.Email, user.Password, time.Now().Format("2006-01-02 15:04:05")).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
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
