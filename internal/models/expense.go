package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Expense struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID   primitive.ObjectID `json:"userId" bson:"userId"`
	Icon     string             `json:"icon" bson:"icon"`
	Category string             `json:"category" bson:"category"`
	Amount   float64            `json:"amount" bson:"amount"`
	Date     time.Time          `json:"date" bson:"date"`
}
