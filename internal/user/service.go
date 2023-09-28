package user

type service struct {
	repo Repository
}

// returns service struct implementing the Service interface
func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) CreateUser(user *User) (int64, error) {
	id, err := s.repo.CreateUser(user)
	if err != nil {
		return 0, nil
	}

	return id, nil
}

func (s *service) GetUserByEmail(email string, password string) (user *User, err error) {
	return
}
