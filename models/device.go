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

type DeviceStatus string

const (
	InStock            DeviceStatus = "InStock"
	AssignedToCustomer DeviceStatus = "AssignedToCustomer"
	AssignedToBuilding DeviceStatus = "AssignedToBuilding"
	PendingCollection  DeviceStatus = "PendingCollection"
)

type Device struct {
	ID             string  `gorm:"primaryKey" json:"id"`
	BuildingID     *string `gorm:"index;foreignKey:ID" json:"buildingId,omitempty"`
	CustomerID     string  `gorm:"index" json:"customerId,omitempty"`
	SubscriptionID string  `gorm:"index" json:"subscriptionId,omitempty"`
	SerialNumber   string  `gorm:"uniqueIndex;type:varchar(50)" json:"serialNumber"`
	Brand          string  `gorm:"type:varchar(50)" json:"brand"`
	Model          string  `gorm:"type:varchar(50)" json:"model"`

	PurchaseDate   *time.Time `json:"purchaseDate,omitempty"`
	AssignedDate   *time.Time `json:"assignedDate,omitempty"`
	CollectionDate *time.Time `json:"collectionDate,omitempty"`

	Type      DeviceType   `gorm:"type:varchar(20)" json:"type"`
	Usage     DeviceUsage  `gorm:"type:varchar(20)" json:"usage"`
	Status    DeviceStatus `gorm:"type:varchar(20)" json:"status"`
	CreatedAt time.Time    `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time    `gorm:"autoUpdateTime" json:"updatedAt"`
}
