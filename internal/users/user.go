package users

import (
	"time"
)

type session struct {
	userId  int64
	expires time.Time
}

func (s *session) isExpired() bool {
	return s.expires.Before(time.Now())
}

var sessions = make(map[string]session)

type User struct {
	ID        int64     `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"name"`
	Password  string    `json:"password" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	LastLogin time.Time `json:"last_login" db:"last_login"`
}

type UserResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Repository interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	SetLastLogin(id int64, lastLogin time.Time) error
}

type Service interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	SetLastLogin(id int64, lastLogin time.Time) error
}

func IsAuthenticated(sessionToken string) bool {
	userSession, exists := sessions[sessionToken]
	if !exists {
		return false
	}

	if userSession.isExpired() {
		delete(sessions, sessionToken)
		return false
	}

	return true
}

func GetIdBySession(sessionToken string) int64 {
	return sessions[sessionToken].userId
}
