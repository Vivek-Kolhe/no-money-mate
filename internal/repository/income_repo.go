package repository

import (
	"context"
	"time"

	_err "github.com/Vivek-Kolhe/no-money-mate/internal/errors"
	"github.com/Vivek-Kolhe/no-money-mate/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IncomeRepository struct {
	db *models.Database
}

func NewIncomeRepository(db *models.Database) *IncomeRepository {
	return &IncomeRepository{db: db}
}

func (r *IncomeRepository) CreateIncome(income models.Income) (primitive.ObjectID, error) {
	collection := r.db.GetCollection("incomes")
	res, err := collection.InsertOne(context.TODO(), income)
	if err != nil {
		return primitive.NilObjectID, err
	}

	insertedID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, _err.ErrFailedToConvertInsertedIdToObjectId
	}

	return insertedID, nil
}

func (r *IncomeRepository) GetIncomeByUserID(userID primitive.ObjectID, month time.Month, year int) ([]models.Income, error) {
	incomes := make([]models.Income, 0)
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	filter := bson.M{
		"userId": userID,
		"date": bson.M{
			"$gte": startDate,
			"$lt":  endDate,
		},
	}

	cursor, err := r.db.GetCollection("incomes").Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var income models.Income

		if err := cursor.Decode(&income); err != nil {
			return nil, err
		}

		incomes = append(incomes, income)
	}

	return incomes, nil
}

func (r *IncomeRepository) DeleteIncome(incomeID, userID primitive.ObjectID) error {
	collection := r.db.GetCollection("incomes")

	filter := bson.M{
		"_id":    incomeID,
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
