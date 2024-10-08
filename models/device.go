package models

import "time"

type DeviceType string
type DeviceUsage string

const (
	ONU    DeviceType = "ONU"
	Switch DeviceType = "Switch"
	OLT    DeviceType = "OLT"
)

const (
	CustomerUse DeviceUsage = "CustomerUse"
	BuildingUse DeviceUsage = "BuildingUse"
	CompanyUse  DeviceUsage = "CompanyUse"
)

type Device struct {
	ID           string      `gorm:"primaryKey" json:"id"`
	Type         DeviceType  `gorm:"type:varchar(20)" json:"type"`
	SerialNumber string      `gorm:"uniqueIndex;type:varchar(50)" json:"serialNumber"`
	PurchaseDate time.Time   `json:"purchaseDate"`
	Usage        DeviceUsage `gorm:"type:varchar(20)" json:"usage"`
	CustomerID   string      `gorm:"index;foreignKey:ID" json:"customerId,omitempty"`
	BuildingID   string      `gorm:"index;foreignKey:ID" json:"buildingId,omitempty"`
	Status       string      `gorm:"type:varchar(20)" json:"status"`
	AssignedDate time.Time   `json:"assignedDate,omitempty"`
	ReturnDate   time.Time   `json:"returnDate,omitempty"`
	CreatedAt    time.Time   `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time   `gorm:"autoUpdateTime" json:"updatedAt"`
}
