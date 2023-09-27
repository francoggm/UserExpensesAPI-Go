package user

type service struct {
	repo *repository
}

func NewService(r *repository) *service {
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
