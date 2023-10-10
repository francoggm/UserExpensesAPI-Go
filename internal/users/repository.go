package users

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

	createdDate := time.Now()

	err := r.db.QueryRow(query, user.Email, user.Password, user.Name, createdDate.Format("2006-01-02 15:04:05")).Scan(&id)
	if err != nil {
		return err
	}

	user.ID = id
	user.CreatedAt = createdDate

	return nil
}

func (r *repository) GetUserByEmail(email string) (*User, error) {
	var user User
	var lastLogin sql.NullTime

	err := r.db.QueryRow("SELECT * FROM users WHERE email=$1", email).Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.CreatedAt, &lastLogin)
	if err != nil {
		return nil, err
	}

	user.LastLogin = lastLogin.Time

	return &user, nil
}

func (r *repository) SetLastLogin(id int64, lastLogin time.Time) error {
	query := "UPDATE users SET last_login=$1 WHERE id=$2"

	_, err := r.db.Exec(query, lastLogin.Format("2006-01-02 15:04:05"), id)
	return err
}
