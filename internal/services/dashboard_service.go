package services

import (
	"github.com/Vivek-Kolhe/no-money-mate/internal/models"
	"github.com/Vivek-Kolhe/no-money-mate/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DashboardService struct {
	repo *repository.DashboardRepository
}

func NewDashboardService(repo *repository.DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetDashboardData(userID primitive.ObjectID) (*models.DashboardData, error) {
	return s.repo.GetDashboardData(userID)
}
