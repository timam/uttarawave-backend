package response

import (
	"github.com/timam/uttarawave-backend/models"
	"time"
)

type DeviceResponse struct {
	ID           string              `json:"id"`
	Type         models.DeviceType   `json:"type"`
	Status       models.DeviceStatus `json:"status"`
	Usage        models.DeviceUsage  `json:"usage"`
	SerialNumber string              `json:"serialNumber"`
	Brand        string              `json:"brand"`
	Model        string              `json:"model"`
	PurchaseDate time.Time           `json:"purchaseDate"`
	AssignedTo   string              `json:"assignedTo,omitempty"`
	AssignedDate *time.Time          `json:"assignedDate,omitempty"`
	CreatedAt    time.Time           `json:"createdAt"`
	UpdatedAt    time.Time           `json:"updatedAt"`
}

func NewDeviceResponse(device *models.Device) DeviceResponse {
	response := DeviceResponse{
		ID:           device.ID,
		Type:         device.Type,
		Status:       device.Status,
		Usage:        device.Usage,
		SerialNumber: device.SerialNumber,
		Brand:        device.Brand,
		Model:        device.Model,
		PurchaseDate: *device.PurchaseDate,
		CreatedAt:    device.CreatedAt,
		UpdatedAt:    device.UpdatedAt,
	}

	if device.SubscriptionID != nil {
		response.AssignedTo = "Subscription"
		response.AssignedDate = device.AssignedDate
	} else if device.BuildingID != nil {
		response.AssignedTo = "Building"
		response.AssignedDate = device.AssignedDate
	}

	return response
}

type DeviceListResponse struct {
	Devices []DeviceResponse `json:"devices"`
	Total   int64            `json:"total"`
	Page    int              `json:"page"`
	Size    int              `json:"size"`
}

func NewDeviceListResponse(devices []models.Device, total int64, page, size int) DeviceListResponse {
	deviceResponses := make([]DeviceResponse, len(devices))
	for i, device := range devices {
		deviceResponses[i] = NewDeviceResponse(&device)
	}

	return DeviceListResponse{
		Devices: deviceResponses,
		Total:   total,
		Page:    page,
		Size:    size,
	}
}
