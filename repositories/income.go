package repositories

import (
	"context"
	"github.com/timam/uttarawave-backend/models"
	"github.com/timam/uttarawave-backend/pkg/db"
)

type IncomeRepository interface {
	CreateIncome(ctx context.Context, income *models.Income) error
	GetIncomeByID(ctx context.Context, id string) (*models.Income, error)
	GetAllIncomes(ctx context.Context) ([]models.Income, error)
	UpdateIncome(ctx context.Context, income *models.Income) error
	DeleteIncome(ctx context.Context, id string) error
}

type GormIncomeRepository struct{}

func NewGormIncomeRepository() *GormIncomeRepository {
	return &GormIncomeRepository{}
}

func (r *GormIncomeRepository) CreateIncome(ctx context.Context, income *models.Income) error {
	return db.DB.WithContext(ctx).Create(income).Error
}

func (r *GormIncomeRepository) GetIncomeByID(ctx context.Context, id string) (*models.Income, error) {
	var income models.Income
	err := db.DB.WithContext(ctx).First(&income, "id = ?", id).Error
	return &income, err
}

func (r *GormIncomeRepository) GetAllIncomes(ctx context.Context) ([]models.Income, error) {
	var incomes []models.Income
	err := db.DB.WithContext(ctx).Find(&incomes).Error
	return incomes, err
}

func (r *GormIncomeRepository) UpdateIncome(ctx context.Context, income *models.Income) error {
	return db.DB.WithContext(ctx).Save(income).Error
}

func (r *GormIncomeRepository) DeleteIncome(ctx context.Context, id string) error {
	return db.DB.WithContext(ctx).Delete(&models.Income{}, "id = ?", id).Error
}
