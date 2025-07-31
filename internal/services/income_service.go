package services

import (
	"time"

	"github.com/Vivek-Kolhe/no-money-mate/internal/models"
	"github.com/Vivek-Kolhe/no-money-mate/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IncomeService struct {
	repo *repository.IncomeRepository
}

func NewIncomeService(repo *repository.IncomeRepository) *IncomeService {
	return &IncomeService{repo: repo}
}

func (s *IncomeService) CreateIncome(income models.Income) (primitive.ObjectID, error) {
	return s.repo.CreateIncome(income)
}

func (s *IncomeService) GetIncomeByUserID(userID primitive.ObjectID, month time.Month, year int) ([]models.Income, error) {
	return s.repo.GetIncomeByUserID(userID, month, year)
}

func (s *IncomeService) DeleteIncome(incomeID, userID primitive.ObjectID) error {
	return s.repo.DeleteIncome(incomeID, userID)
}
