package response

import (
	"github.com/timam/uttarawave-backend/internals/models"
	"time"
)

type DeviceResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type DeviceListResponse struct {
	Items      []DeviceItemResponse `json:"items"`
	Pagination PaginationInfo       `json:"pagination"`
}

type DeviceItemResponse struct {
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

func NewDeviceResponse(status int, message string, data interface{}) DeviceResponse {
	return DeviceResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func NewDeviceItemResponse(device *models.Device) DeviceItemResponse {
	response := DeviceItemResponse{
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

func NewDeviceListResponse(devices []models.Device, total int64, page, size int) DeviceListResponse {
	deviceResponses := make([]DeviceItemResponse, len(devices))
	for i, device := range devices {
		deviceResponses[i] = NewDeviceItemResponse(&device)
	}

	return DeviceListResponse{
		Items: deviceResponses,
		Pagination: PaginationInfo{
			Total: total,
			Page:  page,
			Size:  size,
		},
	}
}
