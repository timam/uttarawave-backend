package repositories

import (
	"context"
	"github.com/timam/uttarawave-backend/models"
	"github.com/timam/uttarawave-backend/pkg/db"
)

type ExpenseRepository interface {
	CreateExpense(ctx context.Context, expense *models.Expense) error
	GetExpenseByID(ctx context.Context, id string) (*models.Expense, error)
	GetAllExpenses(ctx context.Context) ([]models.Expense, error)
	UpdateExpense(ctx context.Context, expense *models.Expense) error
	DeleteExpense(ctx context.Context, id string) error
}

type GormExpenseRepository struct{}

func NewGormExpenseRepository() *GormExpenseRepository {
	return &GormExpenseRepository{}
}

func (r *GormExpenseRepository) CreateExpense(ctx context.Context, expense *models.Expense) error {
	return db.DB.WithContext(ctx).Create(expense).Error
}

func (r *GormExpenseRepository) GetExpenseByID(ctx context.Context, id string) (*models.Expense, error) {
	var expense models.Expense
	err := db.DB.WithContext(ctx).First(&expense, "id = ?", id).Error
	return &expense, err
}

func (r *GormExpenseRepository) GetAllExpenses(ctx context.Context) ([]models.Expense, error) {
	var expenses []models.Expense
	err := db.DB.WithContext(ctx).Find(&expenses).Error
	return expenses, err
}

func (r *GormExpenseRepository) UpdateExpense(ctx context.Context, expense *models.Expense) error {
	return db.DB.WithContext(ctx).Save(expense).Error
}

func (r *GormExpenseRepository) DeleteExpense(ctx context.Context, id string) error {
	return db.DB.WithContext(ctx).Delete(&models.Expense{}, "id = ?", id).Error
}
