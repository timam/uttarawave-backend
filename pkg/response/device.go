package response

import "time"

// DeviceResponse represents a single device in API responses
type DeviceResponse struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Status       string    `json:"status"`
	SerialNumber string    `json:"serialNumber"`
	Brand        string    `json:"brand"`
	Model        string    `json:"model"`
	PurchaseDate time.Time `json:"purchaseDate"`
	AssignedTo   string    `json:"assignedTo,omitempty"`
	AssignedDate time.Time `json:"assignedDate,omitempty"`
}

// DeviceListResponse represents a paginated list of devices
type DeviceListResponse struct {
	Devices     []DeviceResponse `json:"devices"`
	TotalCount  int64            `json:"totalCount"`
	CurrentPage int              `json:"currentPage"`
	PageSize    int              `json:"pageSize"`
}

// DeviceAssignmentResponse represents the result of a device assignment operation
type DeviceAssignmentResponse struct {
	DeviceID       string    `json:"deviceId"`
	AssignmentType string    `json:"assignmentType"`
	AssignmentID   string    `json:"assignmentId"`
	AssignedDate   time.Time `json:"assignedDate"`
}
