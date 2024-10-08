package repositories

import (
	"context"
	"github.com/timam/uttarawave-backend/models"
	"github.com/timam/uttarawave-backend/pkg/db"
	"time"
)

type DeviceRepository interface {
	CreateDevice(ctx context.Context, device *models.Device) error
	GetDeviceByID(ctx context.Context, id string) (*models.Device, error)
	GetAllDevices(ctx context.Context) ([]models.Device, error)
	UpdateDevice(ctx context.Context, device *models.Device) error
	DeleteDevice(ctx context.Context, id string) error
	AssignDeviceToCustomer(ctx context.Context, deviceID, customerID string) error
	AssignDeviceToBuilding(ctx context.Context, deviceID, buildingID string) error
	UnassignDevice(ctx context.Context, deviceID string) error
}

type GormDeviceRepository struct{}

func NewGormDeviceRepository() *GormDeviceRepository {
	return &GormDeviceRepository{}
}

func (r *GormDeviceRepository) CreateDevice(ctx context.Context, device *models.Device) error {
	return db.DB.WithContext(ctx).Create(device).Error
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

func (r *GormDeviceRepository) AssignDeviceToCustomer(ctx context.Context, deviceID, customerID string) error {
	return db.DB.WithContext(ctx).Model(&models.Device{}).Where("id = ?", deviceID).Updates(map[string]interface{}{
		"customer_id":   customerID,
		"building_id":   nil,
		"usage":         models.CustomerUse,
		"assigned_date": time.Now(),
		"return_date":   nil,
	}).Error
}

func (r *GormDeviceRepository) AssignDeviceToBuilding(ctx context.Context, deviceID, buildingID string) error {
	return db.DB.WithContext(ctx).Model(&models.Device{}).Where("id = ?", deviceID).Updates(map[string]interface{}{
		"building_id":   buildingID,
		"customer_id":   nil,
		"usage":         models.BuildingUse,
		"assigned_date": time.Now(),
		"return_date":   nil,
	}).Error
}

func (r *GormDeviceRepository) UnassignDevice(ctx context.Context, deviceID string) error {
	return db.DB.WithContext(ctx).Model(&models.Device{}).Where("id = ?", deviceID).Updates(map[string]interface{}{
		"customer_id":   nil,
		"building_id":   nil,
		"usage":         models.CompanyUse,
		"return_date":   time.Now(),
		"assigned_date": nil,
	}).Error
}
