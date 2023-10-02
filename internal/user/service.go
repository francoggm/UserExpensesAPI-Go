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

func (s *service) CreateUser(user *User) (error) {
	err := s.repo.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetUserByEmail(email string) (*User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
