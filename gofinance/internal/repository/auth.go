package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yourusername/gofinance/internal/model"
)

type authRepo struct {
	db *sqlx.DB
}

func newAuthRepo(db *sqlx.DB) *authRepo {
	return &authRepo{db: db}
}

func (r *authRepo) CreateUser(email, passwordHash string) (int, error) {
	var id int
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRow(query, email, passwordHash).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}
	return id, nil
}

func (r *authRepo) GetUserByEmail(email string) (*model.User, error) {
	var u model.User
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`
	err := r.db.Get(&u, query, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return &u, nil
}
