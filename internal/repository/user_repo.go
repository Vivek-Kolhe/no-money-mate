package repository

import (
	"context"

	"github.com/Vivek-Kolhe/no-money-mate/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository struct {
	db *models.Database
}

func NewUserRepository(db *models.Database) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user models.User) error {
	collection := r.db.GetCollection("users")
	_, err := collection.InsertOne(context.TODO(), user)
	return err
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	collection := r.db.GetCollection(("users"))

	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
