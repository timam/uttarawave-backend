package repositories

import (
	"context"
	"github.com/timam/uttarawave-backend/internals/models"
	"github.com/timam/uttarawave-backend/pkg/db"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *models.Payment) error
	GetPaymentByID(ctx context.Context, id string) (*models.Payment, error)
	GetAllPayments(ctx context.Context) ([]models.Payment, error)
	UpdatePayment(ctx context.Context, payment *models.Payment) error
	DeletePayment(ctx context.Context, id string) error
}

type GormPaymentRepository struct{}

func NewGormPaymentRepository() *GormPaymentRepository {
	return &GormPaymentRepository{}
}

func (r *GormPaymentRepository) CreatePayment(ctx context.Context, payment *models.Payment) error {
	return db.DB.WithContext(ctx).Create(payment).Error
}

func (r *GormPaymentRepository) GetPaymentByID(ctx context.Context, id string) (*models.Payment, error) {
	var payment models.Payment
	err := db.DB.WithContext(ctx).First(&payment, "id = ?", id).Error
	return &payment, err
}

func (r *GormPaymentRepository) GetAllPayments(ctx context.Context) ([]models.Payment, error) {
	var payments []models.Payment
	err := db.DB.WithContext(ctx).Find(&payments).Error
	return payments, err
}

func (r *GormPaymentRepository) UpdatePayment(ctx context.Context, payment *models.Payment) error {
	return db.DB.WithContext(ctx).Save(payment).Error
}

func (r *GormPaymentRepository) DeletePayment(ctx context.Context, id string) error {
	return db.DB.WithContext(ctx).Delete(&models.Payment{}, "id = ?", id).Error
}
