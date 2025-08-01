package repository

import (
	"context"
	"sort"
	"time"

	"github.com/Vivek-Kolhe/no-money-mate/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type DashboardRepository struct {
	db *models.Database
}

func NewDashboardRepository(db *models.Database) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) GetDashboardData(userID primitive.ObjectID) (*models.DashboardData, error) {
	ctx := context.TODO()
	now := time.Now()
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	last60, last30 := now.AddDate(0, 0, -60), now.AddDate(0, 0, -30)

	var incomeTotalMonth, incomeTotal60 float64
	var expenseTotalMonth, expenseTotal30 float64
	var last30Expenses []models.Expense
	var last60Incomes []models.Income
	var last5Transactions []models.Transaction

	incomeCollection := r.db.GetCollection("incomes")
	expenseCollection := r.db.GetCollection("expenses")

	// Aggregation
	filter := bson.M{"userId": userID, "data": bson.M{"$gte": firstOfMonth}}
	cursor, err := incomeCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var i models.Income
		if err := cursor.Decode(&i); err != nil {
			return nil, err
		}

		incomeTotalMonth += i.Amount
	}
	cursor.Close(ctx)

	filter = bson.M{"userId": userID, "date": bson.M{"$gte": last60}}
	cursor, err = incomeCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var i models.Income
		if err := cursor.Decode(&i); err != nil {
			return nil, err
		}

		incomeTotal60 += i.Amount
		last60Incomes = append(last60Incomes, i)
	}
	cursor.Close(ctx)

	filter = bson.M{"userId": userID, "date": bson.M{"$gte": firstOfMonth}}
	cursor, err = expenseCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var e models.Expense
		if err := cursor.Decode(&e); err != nil {
			return nil, err
		}

		expenseTotalMonth += e.Amount
	}
	cursor.Close(ctx)

	filter = bson.M{"userId": userID, "date": bson.M{"$gte": last30}}
	cursor, err = expenseCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var e models.Expense
		if err := cursor.Decode(&e); err != nil {
			return nil, err
		}

		expenseTotal30 += e.Amount
		last30Expenses = append(last30Expenses, e)
	}
	cursor.Close(ctx)

	// Last 5 transactions
	expenseCursor, err := expenseCollection.Find(ctx, bson.M{"userId": userID}, &options.FindOptions{
		Sort:  bson.M{"date": -1},
		Limit: pointerInt64(5),
	})
	if err != nil {
		return nil, err
	}

	incomeCursor, err := incomeCollection.Find(ctx, bson.M{"userId": userID}, &options.FindOptions{
		Sort:  bson.M{"date": -1},
		Limit: pointerInt64(5),
	})
	if err != nil {
		return nil, err
	}

	for expenseCursor.Next(ctx) {
		var e models.Expense
		if err := expenseCursor.Decode(&e); err != nil {
			return nil, err
		}

		last5Transactions = append(last5Transactions, models.Transaction{
			Type:     "expense",
			Amount:   e.Amount,
			Date:     e.Date,
			Icon:     e.Icon,
			Category: e.Category,
		})
	}

	for incomeCursor.Next(ctx) {
		var i models.Income
		if err := incomeCursor.Decode(&i); err != nil {
			return nil, err
		}

		last5Transactions = append(last5Transactions, models.Transaction{
			Type:   "income",
			Amount: i.Amount,
			Date:   i.Date,
			Icon:   i.Icon,
			Source: i.Source,
		})
	}

	expenseCursor.Close(ctx)
	incomeCursor.Close(ctx)

	sort.SliceStable(last5Transactions, func(i, j int) bool {
		return last5Transactions[i].Date.After(last5Transactions[j].Date)
	})
	if len(last5Transactions) > 5 {
		last5Transactions = last5Transactions[:5]
	}

	return &models.DashboardData{
		MonthIncomeTotal:   incomeTotalMonth,
		MonthExpenseTotal:  expenseTotalMonth,
		Balance:            max(incomeTotalMonth-expenseTotalMonth, 0),
		Last60DaysIncome:   incomeTotal60,
		Last30DaysExpense:  expenseTotal30,
		Last30DaysExpenses: last30Expenses,
		Last60DaysIncomes:  last60Incomes,
		Last5Transactions:  last5Transactions,
	}, nil
}

func pointerInt64(n int64) *int64 {
	return &n
}
