package repositories

import (
	"context"
	"github.com/timam/uttarawave-backend/internals/models"
	"github.com/timam/uttarawave-backend/pkg/db"
)

type InvoiceRepository interface {
	CreateInvoice(ctx context.Context, invoice *models.Invoice) error
	GetInvoiceByID(ctx context.Context, id string) (*models.Invoice, error)
	GetAllInvoices(ctx context.Context) ([]models.Invoice, error)
	UpdateInvoice(ctx context.Context, invoice *models.Invoice) error
	DeleteInvoice(ctx context.Context, id string) error
	GetInvoicesByCustomerID(ctx context.Context, customerID string) ([]models.Invoice, error)
	GetInvoicesBySubscriptionID(ctx context.Context, subscriptionID string) ([]models.Invoice, error)
}

type GormInvoiceRepository struct{}

func NewGormInvoiceRepository() *GormInvoiceRepository {
	return &GormInvoiceRepository{}
}

func (r *GormInvoiceRepository) CreateInvoice(ctx context.Context, invoice *models.Invoice) error {
	return db.DB.WithContext(ctx).Create(invoice).Error
}

func (r *GormInvoiceRepository) GetInvoiceByID(ctx context.Context, id string) (*models.Invoice, error) {
	var invoice models.Invoice
	err := db.DB.WithContext(ctx).First(&invoice, "id = ?", id).Error
	return &invoice, err
}

func (r *GormInvoiceRepository) GetAllInvoices(ctx context.Context) ([]models.Invoice, error) {
	var invoices []models.Invoice
	err := db.DB.WithContext(ctx).Find(&invoices).Error
	return invoices, err
}

func (r *GormInvoiceRepository) UpdateInvoice(ctx context.Context, invoice *models.Invoice) error {
	return db.DB.WithContext(ctx).Save(invoice).Error
}

func (r *GormInvoiceRepository) DeleteInvoice(ctx context.Context, id string) error {
	return db.DB.WithContext(ctx).Delete(&models.Invoice{}, "id = ?", id).Error
}

func (r *GormInvoiceRepository) GetInvoicesByCustomerID(ctx context.Context, customerID string) ([]models.Invoice, error) {
	var invoices []models.Invoice
	err := db.DB.WithContext(ctx).Where("customer_id = ?", customerID).Find(&invoices).Error
	return invoices, err
}

func (r *GormInvoiceRepository) GetInvoicesBySubscriptionID(ctx context.Context, subscriptionID string) ([]models.Invoice, error) {
	var invoices []models.Invoice
	err := db.DB.WithContext(ctx).Where("subscription_id = ?", subscriptionID).Find(&invoices).Error
	return invoices, err
}
