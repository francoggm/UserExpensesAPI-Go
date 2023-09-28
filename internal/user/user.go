package user

type User struct {
	ID       int64  `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type Repository interface {
	CreateUser(user *User) (int64, error)
	GetUserByEmail(email string) (user *User, err error)
}

type Service interface {
	CreateUser(user *User) (int64, error)
	GetUserByEmail(email string, password string) (user *User, err error)
}
