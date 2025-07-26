package repository

import (
	"context"
	"time"

	"github.com/Vivek-Kolhe/no-money-mate/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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

func (r *ExpenseRepository) GetExpensesByUserID(userID primitive.ObjectID, month time.Month, year int) ([]models.Expense, error) {
	expenses := make([]models.Expense, 0)
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	filter := bson.M{
		"userId": userID,
		"date": bson.M{
			"$gte": startDate,
			"$lt":  endDate,
		},
	}

	cursor, err := r.db.GetCollection("expenses").Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var expense models.Expense

		if err := cursor.Decode(&expense); err != nil {
			return nil, err
		}

		expenses = append(expenses, expense)
	}

	return expenses, nil
}

func (r *ExpenseRepository) DeleteExpense(expenseID, userID primitive.ObjectID) error {
	collection := r.db.GetCollection("expenses")

	filter := bson.M{
		"_id":    expenseID,
		"userId": userID,
	}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
