package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/yourusername/gofinance/internal/model"
)

type Repository struct {
	Auth
	Transaction
	Category
}

type Auth interface {
	CreateUser(email, passwordHash string) (int, error)
	GetUserByEmail(email string) (*model.User, error)
}

type Transaction interface {
	Create(userID int, input model.TransactionInput) (int, error)
	GetAll(userID int) ([]model.Transaction, error)
	GetByID(userID, id int) (*model.Transaction, error)
	Update(userID, id int, input model.TransactionInput) error
	Delete(userID, id int) error
	GetSummary(userID int, from, to string) (*model.Summary, error)
}

type Category interface {
	GetAll() ([]model.Category, error)
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		Auth:        newAuthRepo(db),
		Transaction: newTransactionRepo(db),
		Category:    newCategoryRepo(db),
	}
}
