package repositories

import (
	"context"
	"github.com/timam/uttarawave-backend/models"
	"github.com/timam/uttarawave-backend/pkg/db"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction *models.Transaction) error
	GetTransactionByID(ctx context.Context, id string) (*models.Transaction, error)
	GetTransactionsBySubscriptionID(ctx context.Context, subscriptionID string) ([]models.Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *models.Transaction) error
}

type GormTransactionRepository struct{}

func NewGormTransactionRepository() *GormTransactionRepository {
	return &GormTransactionRepository{}
}

func (r *GormTransactionRepository) CreateTransaction(ctx context.Context, transaction *models.Transaction) error {
	return db.DB.WithContext(ctx).Create(transaction).Error
}

func (r *GormTransactionRepository) GetTransactionByID(ctx context.Context, id string) (*models.Transaction, error) {
	var transaction models.Transaction
	err := db.DB.WithContext(ctx).First(&transaction, "id = ?", id).Error
	return &transaction, err
}

func (r *GormTransactionRepository) GetTransactionsBySubscriptionID(ctx context.Context, subscriptionID string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := db.DB.WithContext(ctx).Where("subscription_id = ?", subscriptionID).Find(&transactions).Error
	return transactions, err
}

func (r *GormTransactionRepository) UpdateTransaction(ctx context.Context, transaction *models.Transaction) error {
	return db.DB.WithContext(ctx).Save(transaction).Error
}
