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
	AssignDevice(ctx context.Context, deviceID string, assignmentType string, assignmentID string) error
	UnassignDevice(ctx context.Context, deviceID string) error
	GetDeviceByAssignment(ctx context.Context, assignmentType string, assignmentID string) (*models.Device, error)
	MarkDeviceStatus(ctx context.Context, deviceID string, status models.DeviceStatus) error
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
	if err != nil {
		return nil, err
	}
	return &device, nil
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

func (r *GormDeviceRepository) AssignDevice(ctx context.Context, deviceID string, assignmentType string, assignmentID string) error {
	updates := map[string]interface{}{
		"status":        models.Assigned,
		"assigned_date": time.Now(),
	}

	switch assignmentType {
	case "Subscription":
		updates["subscription_id"] = assignmentID
		updates["building_id"] = nil
	case "Building":
		updates["building_id"] = assignmentID
		updates["subscription_id"] = nil
	default:
		return errors.New("invalid assignment type")
	}

	return db.DB.WithContext(ctx).Model(&models.Device{}).Where("id = ?", deviceID).Updates(updates).Error
}

func (r *GormDeviceRepository) UnassignDevice(ctx context.Context, deviceID string) error {
	return db.DB.WithContext(ctx).Model(&models.Device{}).Where("id = ?", deviceID).Updates(map[string]interface{}{
		"subscription_id": nil,
		"building_id":     nil,
		"status":          models.InStock,
		"assigned_date":   nil,
		"collection_date": nil,
	}).Error
}

func (r *GormDeviceRepository) GetDeviceByAssignment(ctx context.Context, assignmentType string, assignmentID string) (*models.Device, error) {
	var device models.Device
	var err error

	switch assignmentType {
	case "Subscription":
		err = db.DB.WithContext(ctx).Where("subscription_id = ?", assignmentID).First(&device).Error
	case "Building":
		err = db.DB.WithContext(ctx).Where("building_id = ?", assignmentID).First(&device).Error
	default:
		return nil, errors.New("invalid assignment type")
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &device, nil
}

func (r *GormDeviceRepository) MarkDeviceStatus(ctx context.Context, deviceID string, status models.DeviceStatus) error {
	return db.DB.WithContext(ctx).Model(&models.Device{}).Where("id = ?", deviceID).Update("status", status).Error
}
