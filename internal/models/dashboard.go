package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DashboardData struct {
	MonthIncomeTotal   float64       `json:"monthIncomeTotal"`
	MonthExpenseTotal  float64       `json:"monthExpenseTotal"`
	Balance            float64       `json:"balance"`
	Last60DaysIncome   float64       `json:"last60DaysIncome"`
	Last30DaysExpense  float64       `json:"last30DaysExpense"`
	Last30DaysExpenses []Expense     `json:"last30DaysExpenses"`
	Last60DaysIncomes  []Income      `json:"last60DaysIncomes"`
	Last5Transactions  []Transaction `json:"last5Transactions"`
}

type Transaction struct {
	ID       primitive.ObjectID `json:"id,omitempty"`
	Type     string             `json:"type"`
	Amount   float64            `json:"amount"`
	Date     time.Time          `json:"date"`
	Icon     string             `json:"icon"`
	Source   string             `json:"source"`
	Category string             `json:"category"`
}
