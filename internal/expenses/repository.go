package expenses

import (
	"database/sql"
	"time"
)

type repository struct {
	db *sql.DB
}

// returns repository struct implementing the Repository interface
func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetExpenses(userId int64) ([]*Expense, error) {
	var expenses []*Expense

	rows, err := r.db.Query("SELECT * FROM expenses WHERE user_id=$1", userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var expense Expense
		var description sql.NullString
		var categoryType, movimentationType sql.NullInt16

		err = rows.Scan(
			&expense.ID,
			&expense.UserID,
			&expense.Title,
			&description,
			&expense.CreatedAt,
			&expense.Value,
			&categoryType,
			&movimentationType,
		)

		if err != nil {
			continue
		}

		expense.Description = description.String
		expense.Category = CategoryType(categoryType.Int16)
		expense.Movimentation = MovimentationType(movimentationType.Int16)

		expenses = append(expenses, &expense)
	}

	return expenses, nil
}

func (r *repository) GetExpense(id, userId int64) (*Expense, error) {
	var expense Expense
	var description sql.NullString
	var categoryType, movimentationType sql.NullInt16

	err := r.db.QueryRow("SELECT * FROM users WHERE id=$1 AND user_id=$2", id, userId).Scan(
		&expense.ID,
		&expense.UserID,
		&expense.Title,
		&description,
		&expense.CreatedAt,
		&expense.Value,
		&categoryType,
		&movimentationType,
	)

	if err != nil {
		return nil, err
	}

	expense.Description = description.String
	expense.Category = CategoryType(categoryType.Int16)
	expense.Movimentation = MovimentationType(movimentationType.Int16)

	return &expense, nil
}

func (r *repository) CreateExpense(expense *Expense) error {
	var id int64

	query := "INSERT INTO expenses (user_id, title, description, value, category_type, movimentation_type, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"

	createdDate := time.Now().Format("2006-01-02 15:04:05")

	err := r.db.QueryRow(query, expense.UserID, expense.Title, expense.Description, expense.Value, expense.Category, expense.Movimentation, createdDate).Scan(&id)
	if err != nil {
		return err
	}

	expense.ID = id
	expense.CreatedAt = createdDate

	return nil
}

func (r *repository) UpdateExpense(expense *Expense) error {
	return nil
}

func (r *repository) DeleteExpense(expense *Expense) error {
	return nil
}
