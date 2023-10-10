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

func (r *repository) ListExpenses(userId int64) ([]*Expense, error) {
	var expenses []*Expense

	query := "SELECT * FROM expenses WHERE user_id=$1"

	rows, err := r.db.Query(query, userId)
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

	query := "SELECT * FROM expenses WHERE id=$1 AND user_id=$2"

	err := r.db.QueryRow(query, id, userId).Scan(
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

	createdDate := time.Now()

	err := r.db.QueryRow(query, expense.UserID, expense.Title, expense.Description, expense.Value, expense.Category, expense.Movimentation, createdDate.Format("2006-01-02 15:04:05")).Scan(&id)
	if err != nil {
		return err
	}

	expense.ID = id
	expense.CreatedAt = createdDate

	return nil
}

func (r *repository) UpdateExpense(id int64, userId int64, expense *Expense) error {
	query := "UPDATE expenses SET title=$1, description=$2, value=$3, category_type=$4, movimentation_type=$5 WHERE id=$6 AND user_id=$7"

	_, err := r.db.Exec(query, expense.Title, expense.Description, expense.Value, expense.Category, expense.Movimentation, id, userId)
	return err
}

func (r *repository) DeleteExpense(id, userId int64) error {
	query := "DELETE FROM expenses WHERE id=$1 AND user_id=$2"

	_, err := r.db.Exec(query, id, userId)
	return err
}
