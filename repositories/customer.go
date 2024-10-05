package repositories

import (
	"context"
	"github.com/timam/uttarawave-backend/models"
	"github.com/timam/uttarawave-backend/pkg/db"
)

type CustomerRepository interface {
	CreateCustomer(ctx context.Context, customer *models.Customer) error
	GetCustomer(id string) (*models.Customer, error)
	GetCustomerByMobile(mobile string) (*models.Customer, error)
	GetAllCustomers() ([]models.Customer, error)
	UpdateCustomer(customer *models.Customer) error
	DeleteCustomer(id string) error
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
	err := db.DB.First(&customer, "id = ?", id).Error
	return &customer, err
}

func (r *GormCustomerRepository) GetCustomerByMobile(mobile string) (*models.Customer, error) {
	var customer models.Customer
	err := db.DB.First(&customer, "mobile = ?", mobile).Error
	return &customer, err
}

func (r *GormCustomerRepository) GetAllCustomers() ([]models.Customer, error) {
	var customers []models.Customer
	err := db.DB.Find(&customers).Error
	return customers, err
}

func (r *GormCustomerRepository) UpdateCustomer(customer *models.Customer) error {
	return db.DB.Save(customer).Error
}

func (r *GormCustomerRepository) DeleteCustomer(id string) error {
	return db.DB.Delete(&models.Customer{}, "id = ?", id).Error
}
