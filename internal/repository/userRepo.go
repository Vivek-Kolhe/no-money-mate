package repository

import (
	"context"

	"github.com/Vivek-Kolhe/no-money-mate/internal/models"
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
