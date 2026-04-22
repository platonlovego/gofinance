package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yourusername/gofinance/internal/model"
)

type transactionRepo struct {
	db *sqlx.DB
}

func newTransactionRepo(db *sqlx.DB) *transactionRepo {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) Create(userID int, input model.TransactionInput) (int, error) {
	var id int
	query := `
		INSERT INTO transactions (user_id, category_id, amount, comment, date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	err := r.db.QueryRow(query, userID, input.CategoryID, input.Amount, input.Comment, input.Date).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create transaction: %w", err)
	}
	return id, nil
}

func (r *transactionRepo) GetAll(userID int) ([]model.Transaction, error) {
	var list []model.Transaction
	query := `
		SELECT id, user_id, category_id, amount, comment, date, created_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY date DESC`
	if err := r.db.Select(&list, query, userID); err != nil {
		return nil, fmt.Errorf("get transactions: %w", err)
	}
	return list, nil
}

func (r *transactionRepo) GetByID(userID, id int) (*model.Transaction, error) {
	var t model.Transaction
	query := `
		SELECT id, user_id, category_id, amount, comment, date, created_at
		FROM transactions
		WHERE id = $1 AND user_id = $2`
	if err := r.db.Get(&t, query, id, userID); err != nil {
		return nil, fmt.Errorf("get transaction by id: %w", err)
	}
	return &t, nil
}

func (r *transactionRepo) Update(userID, id int, input model.TransactionInput) error {
	query := `
		UPDATE transactions
		SET category_id = $1, amount = $2, comment = $3, date = $4
		WHERE id = $5 AND user_id = $6`
	res, err := r.db.Exec(query, input.CategoryID, input.Amount, input.Comment, input.Date, id, userID)
	if err != nil {
		return fmt.Errorf("update transaction: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("transaction not found")
	}
	return nil
}

func (r *transactionRepo) Delete(userID, id int) error {
	res, err := r.db.Exec(`DELETE FROM transactions WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("delete transaction: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("transaction not found")
	}
	return nil
}

func (r *transactionRepo) GetSummary(userID int, from, to string) (*model.Summary, error) {
	var rows []model.CategorySummary
	query := `
		SELECT c.name AS category_name, c.type, COALESCE(SUM(t.amount), 0) AS total
		FROM categories c
		LEFT JOIN transactions t
			ON t.category_id = c.id
			AND t.user_id = $1
			AND t.date BETWEEN $2::date AND $3::date
		GROUP BY c.name, c.type
		ORDER BY c.type, total DESC`
	if err := r.db.Select(&rows, query, userID, from, to); err != nil {
		return nil, fmt.Errorf("get summary: %w", err)
	}

	s := &model.Summary{ByCategory: rows}
	for _, row := range rows {
		if row.Type == "income" {
			s.TotalIncome += row.Total
		} else {
			s.TotalExpense += row.Total
		}
	}
	s.Balance = s.TotalIncome - s.TotalExpense
	return s, nil
}
