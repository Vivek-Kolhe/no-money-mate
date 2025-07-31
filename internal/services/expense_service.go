package services

import (
	"time"

	"github.com/Vivek-Kolhe/no-money-mate/internal/models"
	"github.com/Vivek-Kolhe/no-money-mate/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExpenseService struct {
	repo *repository.ExpenseRepository
}

func NewExpenseService(repo *repository.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: repo}
}

func (s *ExpenseService) CreateExpense(expense models.Expense) (primitive.ObjectID, error) {
	return s.repo.CreateExpense(expense)
}

func (s *ExpenseService) GetExpensesByUserID(userID primitive.ObjectID, month time.Month, year int) ([]models.Expense, error) {
	return s.repo.GetExpensesByUserID(userID, month, year)
}

func (s *ExpenseService) DeleteExpense(expenseId, userId primitive.ObjectID) error {
	return s.repo.DeleteExpense(expenseId, userId)
}
