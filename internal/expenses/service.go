package expenses

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) GetExpenses(userId int64) ([]*ExpenseResponse, error) {
	return s.repo.GetExpenses(userId)
}

func (s *service) GetExpense(id, userId int64) (*ExpenseResponse, error) {
	return s.repo.GetExpense(id, userId)
}

func (s *service) CreateExpense(expense *Expense) error {
	return s.repo.CreateExpense(expense)
}

func (s *service) UpdateExpense(expense *Expense) error {
	return nil
}

func (s *service) DeleteExpense(id, userId int64) error {
	return s.repo.DeleteExpense(id, userId)
}
