package users

type service struct {
	repo Repository
}

// returns service struct implementing the Service interface
func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) CreateUser(user *User) error {
	return s.repo.CreateUser(user)
}

func (s *service) GetUserByEmail(email string) (*User, error) {
	return s.repo.GetUserByEmail(email)
}
