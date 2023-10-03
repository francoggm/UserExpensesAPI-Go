package users

type User struct {
	ID        int64  `json:"id" db:"id"`
	Email     string `json:"email" db:"email"`
	Name      string `json:"name" db:"name"`
	Password  string `json:"password" db:"password"`
	CreatedAt string `json:"created_at" db:"created_at"`
	LastLogin string `json:"last_login" db:"last_login"`
}

type UserResponse struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	LastLogin string `json:"last_login"`
}

type Repository interface {
	CreateUser(user *User) (error)
	GetUserByEmail(email string) (*User, error)
}

type Service interface {
	CreateUser(user *User) (error)
	GetUserByEmail(email string) (*User, error)
}
