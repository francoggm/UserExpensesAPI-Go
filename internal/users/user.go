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
	Email     string    `json:"email" binding:"required,email" db:"email"`
	Name      string    `json:"name" binding:"required" db:"name"`
	Password  string    `json:"password" binding:"required,min=8,max=64" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	LastLogin time.Time `json:"last_login" db:"last_login"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type UserResponse struct {
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	LastLogin time.Time `json:"last_login"`
}

type Repository interface {
	CreateUser(req *RegisterRequest) (*User, error)
	GetUserByEmail(email string) (*User, error)
	SetLastLogin(id int64, lastLogin time.Time) error
}

type Service interface {
	CreateUser(req *RegisterRequest) (*User, error)
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
