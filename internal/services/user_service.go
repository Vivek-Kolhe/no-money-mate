package services

import (
	"errors"

	"github.com/Vivek-Kolhe/no-money-mate/internal/models"
	"github.com/Vivek-Kolhe/no-money-mate/internal/repository"
	"github.com/Vivek-Kolhe/no-money-mate/internal/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user models.User) error {
	_, err := s.repo.FindByEmail(user.Email)
	if err != nil {
		if err.Error() != mongo.ErrNoDocuments.Error() {
			return err
		}
	} else {
		return errors.New("user already exists")
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.repo.CreateUser(user)
}
