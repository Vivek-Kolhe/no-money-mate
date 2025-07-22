package repository

import (
	"context"

	"github.com/Vivek-Kolhe/no-money-mate/internal/models"
)

type ExpenseRepository struct {
	db *models.Database
}

func NewExpenseRepository(db *models.Database) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) CreateExpense(expense models.Expense) error {
	collection := r.db.GetCollection("expenses")
	_, err := collection.InsertOne(context.TODO(), expense)
	return err
}
