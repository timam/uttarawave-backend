package repositories

import (
	"context"
	"errors"
	"github.com/timam/uttarawave-backend/internals/models"
	"github.com/timam/uttarawave-backend/pkg/db"
	"github.com/timam/uttarawave-backend/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	CreateCustomer(ctx context.Context, customer *models.Customer) error
	GetCustomer(id string) (*models.Customer, error)
	GetCustomerByMobile(mobile string) (*models.Customer, error)
	UpdateCustomer(customer *models.Customer) error
	DeleteCustomer(id string) error
	GetCustomersPaginated(page, pageSize int) ([]models.Customer, int64, error)
}

type GormCustomerRepository struct{}

func NewGormCustomerRepository() *GormCustomerRepository {
	return &GormCustomerRepository{}
}

func (r *GormCustomerRepository) CreateCustomer(ctx context.Context, customer *models.Customer) error {
	return db.DB.WithContext(ctx).Create(customer).Error
}

func (r *GormCustomerRepository) GetCustomer(id string) (*models.Customer, error) {
	var customer models.Customer
	result := db.DB.Preload("Address").First(&customer, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.Warn("Customer not found in database", zap.String("id", id))
			return nil, nil
		}
		logger.Error("Database error when fetching customer", zap.Error(result.Error), zap.String("id", id))
		return nil, result.Error
	}
	logger.Info("Customer found in database", zap.String("id", customer.ID), zap.String("name", customer.Name))
	return &customer, nil
}

func (r *GormCustomerRepository) GetCustomerByMobile(mobile string) (*models.Customer, error) {
	var customer models.Customer
	result := db.DB.Preload("Address").Where("mobile = ?", mobile).First(&customer)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.Warn("Customer not found", zap.String("mobile", mobile))
			return nil, nil
		}
		logger.Error("Database error when fetching customer by mobile", zap.Error(result.Error), zap.String("mobile", mobile))
		return nil, result.Error
	}
	logger.Info("Customer found by mobile", zap.String("id", customer.ID), zap.String("mobile", customer.Mobile))
	return &customer, nil
}

func (r *GormCustomerRepository) GetCustomersPaginated(page, pageSize int) ([]models.Customer, int64, error) {
	var customers []models.Customer
	var totalCount int64

	offset := (page - 1) * pageSize

	if err := db.DB.Model(&models.Customer{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if err := db.DB.Preload("Address").Offset(offset).Limit(pageSize).Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	return customers, totalCount, nil
}

func (r *GormCustomerRepository) UpdateCustomer(customer *models.Customer) error {
	return db.DB.Save(customer).Error
}

func (r *GormCustomerRepository) DeleteCustomer(id string) error {
	return db.DB.Delete(&models.Customer{}, "id = ?", id).Error
}
