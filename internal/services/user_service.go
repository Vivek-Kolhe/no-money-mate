package services

import (
	_errors "github.com/Vivek-Kolhe/no-money-mate/internal/errors"
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
		return _errors.ErrUserAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.repo.CreateUser(user)
}

func (s *UserService) GetUser(payload models.LoginRequest) (*models.UserResponseDTO, error) {
	user, err := s.repo.FindByEmail(payload.Email)
	if err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			return nil, _errors.ErrInvalidCredentials
		}
		return nil, err
	}

	isValid := utils.CheckPassword(user.Password, payload.Password)
	if !isValid {
		return nil, _errors.ErrInvalidCredentials
	}

	return &models.UserResponseDTO{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}, nil
}
