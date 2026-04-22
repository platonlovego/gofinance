package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yourusername/gofinance/internal/model"
)

type categoryRepo struct {
	db *sqlx.DB
}

func newCategoryRepo(db *sqlx.DB) *categoryRepo {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) GetAll() ([]model.Category, error) {
	var cats []model.Category
	if err := r.db.Select(&cats, `SELECT id, name, type FROM categories ORDER BY type, name`); err != nil {
		return nil, fmt.Errorf("get categories: %w", err)
	}
	return cats, nil
}
