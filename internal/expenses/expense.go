package expenses

import "time"

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
	CreatedAt     time.Time         `json:"created_at" db:"created_at"`
}

type ExpenseRequest struct {
	Title         string            `json:"title" binding:"required"`
	Description   string            `json:"description"`
	Value         float32           `json:"value" binding:"required"`
	Category      CategoryType      `json:"category_type" binding:"required"`
	Movimentation MovimentationType `json:"movimentation_type" binding:"required"`
}

type ExpenseResponse struct {
	ID            int64             `json:"id" db:"id"`
	Title         string            `json:"title" db:"title"`
	Description   string            `json:"description" db:"description"`
	Value         float32           `json:"value" db:"value"`
	Category      CategoryType      `json:"category_type" db:"category_type"`
	Movimentation MovimentationType `json:"movimentation_type" db:"movimentation_type"`
	CreatedAt     time.Time         `json:"created_at" db:"created_at"`
}

type Repository interface {
	ListExpenses(userId int64) ([]*Expense, error)
	GetExpense(id, userId int64) (*Expense, error)
	CreateExpense(req *ExpenseRequest, userId int64) (*Expense, error)
	UpdateExpense(id int64, userId int64, expense *Expense) error
	DeleteExpense(id, userId int64) error
}

type Service interface {
	ListExpenses(userId int64) ([]*Expense, error)
	GetExpense(id, userId int64) (*Expense, error)
	CreateExpense(req *ExpenseRequest, userId int64) (*Expense, error)
	UpdateExpense(id int64, userId int64, expense *Expense) error
	DeleteExpense(id, userId int64) error
}
