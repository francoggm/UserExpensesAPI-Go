package expenses

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) GetExpenses(userId int64) ([]*Expense, error) {
	return nil, nil
}

func (s *service) GetExpense(id, userId int64) (*Expense, error) {
	return nil, nil
}

func (s *service) CreateExpense(expense *Expense) error {
	return s.repo.CreateExpense(expense)
}

func (s *service) DeleteExpense(expense *Expense) error {
	return nil
}

func (s *service) UpdateExpense(expense *Expense) error {
	return nil
}
