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

func (s *service) Register(email string, password string) (user *User, err error) {
	return
}

func (s *service) Login(email string, password string) (user *User, err error) {
	return
}
