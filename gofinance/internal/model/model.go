package model

import "time"

type User struct {
	ID           int       `db:"id"            json:"id"`
	Email        string    `db:"email"         json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at"    json:"created_at"`
}

type RegisterInput struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Category — справочник категорий (еда, транспорт, зарплата и т.д.)
type Category struct {
	ID   int    `db:"id"   json:"id"`
	Name string `db:"name" json:"name"`
	Type string `db:"type" json:"type"` // "income" | "expense"
}

type Transaction struct {
	ID         int       `db:"id"          json:"id"`
	UserID     int       `db:"user_id"      json:"-"`
	CategoryID int       `db:"category_id"  json:"category_id"`
	Amount     float64   `db:"amount"       json:"amount"`
	Comment    string    `db:"comment"      json:"comment,omitempty"`
	Date       time.Time `db:"date"         json:"date"`
	CreatedAt  time.Time `db:"created_at"   json:"created_at"`
}

type TransactionInput struct {
	CategoryID int       `json:"category_id" binding:"required"`
	Amount     float64   `json:"amount"      binding:"required,gt=0"`
	Comment    string    `json:"comment"`
	Date       time.Time `json:"date"        binding:"required"`
}

// Summary — агрегированная статистика за период
type Summary struct {
	TotalIncome  float64            `json:"total_income"`
	TotalExpense float64            `json:"total_expense"`
	Balance      float64            `json:"balance"`
	ByCategory   []CategorySummary  `json:"by_category"`
}

type CategorySummary struct {
	CategoryName string  `db:"category_name" json:"category_name"`
	Type         string  `db:"type"          json:"type"`
	Total        float64 `db:"total"         json:"total"`
}
