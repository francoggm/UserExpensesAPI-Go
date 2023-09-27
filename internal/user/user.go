package user

type User struct {
	ID int64 `json:"id" db:"id"` 
	Email string `json:"email" db:"email"`
	Password string `json:"password" db:"password"` 
}