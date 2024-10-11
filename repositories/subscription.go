package repositories

import (
	"context"
	"github.com/timam/uttarawave-backend/models"
	"github.com/timam/uttarawave-backend/pkg/db"
	"time"
)

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, subscription *models.Subscription) error
	GetSubscription(ctx context.Context, id string) (*models.Subscription, error)
	GetSubscriptionsByCustomerID(ctx context.Context, customerID string) ([]models.Subscription, error)
	UpdateSubscription(ctx context.Context, subscription *models.Subscription) error
	DeleteSubscription(ctx context.Context, id string) error
	GetSubscriptionsPaginated(ctx context.Context, page, pageSize int) ([]models.Subscription, int64, error)
	GetExpiredSubscriptions(ctx context.Context) ([]models.Subscription, error)
}

type GormSubscriptionRepository struct{}

func NewGormSubscriptionRepository() *GormSubscriptionRepository {
	return &GormSubscriptionRepository{}
}

func (r *GormSubscriptionRepository) CreateSubscription(ctx context.Context, subscription *models.Subscription) error {
	return db.DB.WithContext(ctx).Create(subscription).Error
}

func (r *GormSubscriptionRepository) GetSubscription(ctx context.Context, id string) (*models.Subscription, error) {
	var subscription models.Subscription
	err := db.DB.WithContext(ctx).First(&subscription, "id = ?", id).Error
	return &subscription, err
}

func (r *GormSubscriptionRepository) GetSubscriptionsByCustomerID(ctx context.Context, customerID string) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := db.DB.WithContext(ctx).Where("customer_id = ?", customerID).Find(&subscriptions).Error
	return subscriptions, err
}

func (r *GormSubscriptionRepository) UpdateSubscription(ctx context.Context, subscription *models.Subscription) error {
	return db.DB.WithContext(ctx).Save(subscription).Error
}

func (r *GormSubscriptionRepository) DeleteSubscription(ctx context.Context, id string) error {
	return db.DB.WithContext(ctx).Delete(&models.Subscription{}, "id = ?", id).Error
}

func (r *GormSubscriptionRepository) GetSubscriptionsPaginated(ctx context.Context, page, pageSize int) ([]models.Subscription, int64, error) {
	var subscriptions []models.Subscription
	var totalCount int64

	offset := (page - 1) * pageSize

	if err := db.DB.WithContext(ctx).Model(&models.Subscription{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if err := db.DB.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&subscriptions).Error; err != nil {
		return nil, 0, err
	}

	return subscriptions, totalCount, nil
}
func (r *GormSubscriptionRepository) GetExpiredSubscriptions(ctx context.Context) ([]models.Subscription, error) {
	var expiredSubscriptions []models.Subscription
	err := db.DB.WithContext(ctx).Where("renewal_date <= ? AND status != ?", time.Now(), "Expired").Find(&expiredSubscriptions).Error
	return expiredSubscriptions, err
}
