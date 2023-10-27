package expenses

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) ListExpenses(userId int64) ([]*Expense, error) {
	return s.repo.ListExpenses(userId)
}

func (s *service) GetExpense(id, userId int64) (*Expense, error) {
	return s.repo.GetExpense(id, userId)
}

func (s *service) CreateExpense(req *ExpenseRequest, userId int64) (*Expense, error) {
	return s.repo.CreateExpense(req, userId)
}

func (s *service) UpdateExpense(id int64, userId int64, expense *Expense) error {
	return s.repo.UpdateExpense(id, userId, expense)
}

func (s *service) DeleteExpense(id, userId int64) error {
	return s.repo.DeleteExpense(id, userId)
}
