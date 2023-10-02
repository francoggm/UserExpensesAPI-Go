package user

type User struct {
	ID        int64  `json:"id" db:"user_id"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"password" db:"password"`
	CreatedAt string `json:"created_at" db:"created_at"`
	LastLogin string `json:"last_login" db:"last_login"`
}

type Repository interface {
	CreateUser(user *User) (int64, error)
	GetUserByEmail(email string) (*User, error)
}

type Service interface {
	CreateUser(user *User) (int64, error)
	GetUserByEmail(email string) (*User, error)
}
