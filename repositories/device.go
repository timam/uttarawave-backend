package repositories

import (
	"context"
	"errors"
	"github.com/timam/uttarawave-backend/models"
	"github.com/timam/uttarawave-backend/pkg/db"
	"gorm.io/gorm"
	"time"
)

type DeviceRepository interface {
	CreateDevice(ctx context.Context, device *models.Device) error
	GetDeviceByID(ctx context.Context, id string) (*models.Device, error)
	GetAllDevices(ctx context.Context) ([]models.Device, error)
	UpdateDevice(ctx context.Context, device *models.Device) error
	DeleteDevice(ctx context.Context, id string) error
	AssignDeviceToSubscription(ctx context.Context, deviceID, subscriptionID string) error
	AssignDeviceToBuilding(ctx context.Context, deviceID, buildingID string) error
	UnassignDevice(ctx context.Context, deviceID string) error
	MarkDeviceForCollection(ctx context.Context, deviceID string) error
	ReturnDeviceToStock(ctx context.Context, deviceID string) error
	GetDevicesByStatus(ctx context.Context, status models.DeviceStatus) ([]models.Device, error)
	GetDeviceBySubscriptionID(ctx context.Context, subscriptionID string) (*models.Device, error)
}

type GormDeviceRepository struct{}

func NewGormDeviceRepository() *GormDeviceRepository {
	return &GormDeviceRepository{}
}

func (r *GormDeviceRepository) CreateDevice(ctx context.Context, device *models.Device) error {
	result := db.DB.WithContext(ctx).Create(device)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (r *GormDeviceRepository) GetDeviceByID(ctx context.Context, id string) (*models.Device, error) {
	var device models.Device
	err := db.DB.WithContext(ctx).First(&device, "id = ?", id).Error
	return &device, err
}

func (r *GormDeviceRepository) GetAllDevices(ctx context.Context) ([]models.Device, error) {
	var devices []models.Device
	err := db.DB.WithContext(ctx).Find(&devices).Error
	return devices, err
}

func (r *GormDeviceRepository) UpdateDevice(ctx context.Context, device *models.Device) error {
	return db.DB.WithContext(ctx).Save(device).Error
}

func (r *GormDeviceRepository) DeleteDevice(ctx context.Context, id string) error {
	return db.DB.WithContext(ctx).Delete(&models.Device{}, "id = ?", id).Error
}

func (r *GormDeviceRepository) AssignDeviceToSubscription(ctx context.Context, deviceID, subscriptionID string) error {
	return db.DB.WithContext(ctx).Model(&models.Device{}).Where("id = ?", deviceID).Updates(map[string]interface{}{
		"subscription_id": subscriptionID,
		"status":          models.AssignedToCustomer,
		"assigned_date":   time.Now(),
		"collection_date": nil,
	}).Error
}

func (r *GormDeviceRepository) AssignDeviceToBuilding(ctx context.Context, deviceID, buildingID string) error {
	return db.DB.WithContext(ctx).Model(&models.Device{}).Where("id = ?", deviceID).Updates(map[string]interface{}{
		"building_id":     buildingID,
		"status":          models.AssignedToBuilding,
		"assigned_date":   time.Now(),
		"collection_date": nil,
	}).Error
}

func (r *GormDeviceRepository) UnassignDevice(ctx context.Context, deviceID string) error {
	return db.DB.WithContext(ctx).Model(&models.Device{}).Where("id = ?", deviceID).Updates(map[string]interface{}{
		"customer_id":     nil,
		"building_id":     nil,
		"subscription_id": nil,
		"status":          models.InStock,
		"assigned_date":   nil,
		"collection_date": nil,
	}).Error
}

func (r *GormDeviceRepository) MarkDeviceForCollection(ctx context.Context, deviceID string) error {
	return db.DB.WithContext(ctx).Model(&models.Device{}).Where("id = ?", deviceID).Updates(map[string]interface{}{
		"status":          models.PendingCollection,
		"collection_date": time.Now(),
	}).Error
}

func (r *GormDeviceRepository) ReturnDeviceToStock(ctx context.Context, deviceID string) error {
	return db.DB.WithContext(ctx).Model(&models.Device{}).Where("id = ?", deviceID).Updates(map[string]interface{}{
		"customer_id":     nil,
		"building_id":     nil,
		"subscription_id": nil,
		"status":          models.InStock,
		"assigned_date":   nil,
		"collection_date": nil,
	}).Error
}

func (r *GormDeviceRepository) GetDevicesByStatus(ctx context.Context, status models.DeviceStatus) ([]models.Device, error) {
	var devices []models.Device
	err := db.DB.WithContext(ctx).Where("status = ?", status).Find(&devices).Error
	return devices, err
}

func (r *GormDeviceRepository) GetDeviceBySubscriptionID(ctx context.Context, subscriptionID string) (*models.Device, error) {
	var device models.Device
	err := db.DB.WithContext(ctx).Where("subscription_id = ?", subscriptionID).First(&device).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &device, nil
}
