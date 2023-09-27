package user

type User struct {
	ID       int64  `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type Repository interface {
	CreateUser(email string, password string) (id int64, err error)
	GetUserByEmail(email string) (user *User, err error)
}

type Service interface {
	Register(email string, password string) (user *User, err error)
	Login(email string, password string) (user *User, err error)
}
