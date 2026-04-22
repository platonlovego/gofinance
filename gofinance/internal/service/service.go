package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"github.com/yourusername/gofinance/internal/model"
	"github.com/yourusername/gofinance/internal/repository"
)

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

// ── Auth ─────────────────────────────────────────────────────────────────────

func (s *Service) Register(input model.RegisterInput) (int, error) {
	existing, err := s.repo.GetUserByEmail(input.Email)
	if err != nil {
		return 0, err
	}
	if existing != nil {
		return 0, errors.New("email already taken")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("hash password: %w", err)
	}

	return s.repo.CreateUser(input.Email, string(hash))
}

func (s *Service) Login(input model.LoginInput) (string, error) {
	user, err := s.repo.GetUserByEmail(input.Email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	return s.generateToken(user.ID)
}

func (s *Service) generateToken(userID int) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "changeme_in_production"
	}

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (s *Service) ParseToken(tokenStr string) (int, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "changeme_in_production"
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, errors.New("invalid subject")
	}

	return int(sub), nil
}

// ── Transactions ─────────────────────────────────────────────────────────────

func (s *Service) CreateTransaction(userID int, input model.TransactionInput) (int, error) {
	return s.repo.Create(userID, input)
}

func (s *Service) GetTransactions(userID int) ([]model.Transaction, error) {
	return s.repo.GetAll(userID)
}

func (s *Service) GetTransaction(userID, id int) (*model.Transaction, error) {
	return s.repo.GetByID(userID, id)
}

func (s *Service) UpdateTransaction(userID, id int, input model.TransactionInput) error {
	return s.repo.Update(userID, id, input)
}

func (s *Service) DeleteTransaction(userID, id int) error {
	return s.repo.Delete(userID, id)
}

func (s *Service) GetSummary(userID int, from, to string) (*model.Summary, error) {
	return s.repo.GetSummary(userID, from, to)
}

// ── Categories ───────────────────────────────────────────────────────────────

func (s *Service) GetCategories() ([]model.Category, error) {
	return s.repo.GetAll()
}
