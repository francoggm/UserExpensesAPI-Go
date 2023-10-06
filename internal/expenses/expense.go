package expenses

type CategoryType int
type MovimentationType int

const (
	UndefinedCategory CategoryType = iota
	Food
	Health
	Mobility
	Education
)

const (
	UndefinedMovimentation MovimentationType = iota
	Input
	Output
)

type Expense struct {
	ID            int64             `json:"id" db:"id"`
	UserID        int64             `json:"user_id" db:"user_id"`
	Title         string            `json:"title" db:"title"`
	Description   string            `json:"description" db:"description"`
	Value         float32           `json:"value" db:"value"`
	Category      CategoryType      `json:"category_type" db:"category_type"`
	Movimentation MovimentationType `json:"movimentation_type" db:"movimentation_type"`
	CreatedAt     string            `json:"created_at" db:"created_at"`
}

type Repository interface {
	GetExpenses(userId int64) ([]*Expense, error)
	GetExpense(id, userId int64) (*Expense, error)
	CreateExpense(expense *Expense) error
	DeleteExpense(expense *Expense) error
	UpdateExpense(expense *Expense) error
}

type Service interface {
	GetExpenses(userId int64) ([]*Expense, error)
	GetExpense(id, userId int64) (*Expense, error)
	CreateExpense(expense *Expense) error
	DeleteExpense(expense *Expense) error
	UpdateExpense(expense *Expense) error
}
