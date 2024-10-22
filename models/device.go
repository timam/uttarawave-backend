package models

import "time"

type DeviceType string

const (
	ONU    DeviceType = "ONU"
	Switch DeviceType = "Switch"
	OLT    DeviceType = "OLT"
	Router DeviceType = "ROUTER"
	Server DeviceType = "SERVER"
	Camera DeviceType = "CAMERA"
)

type DeviceUsage string

const (
	CustomerUse DeviceUsage = "CustomerUse"
	BuildingUse DeviceUsage = "BuildingUse"
	CompanyUse  DeviceUsage = "CompanyUse"
)

type DeviceStatus string

const (
	InStock           DeviceStatus = "InStock"
	Assigned          DeviceStatus = "Assigned"
	PendingCollection DeviceStatus = "PendingCollection"
	Damaged           DeviceStatus = "Damaged"
	UnderRepair       DeviceStatus = "UnderRepair"
)

type Device struct {
	ID           string       `gorm:"primaryKey" json:"id"`
	Type         DeviceType   `gorm:"type:varchar(20)" json:"type"`
	Status       DeviceStatus `gorm:"type:varchar(20)" json:"status"`
	SerialNumber string       `gorm:"uniqueIndex;type:varchar(50)" json:"serialNumber"`
	Brand        string       `gorm:"type:varchar(50)" json:"brand"`
	Model        string       `gorm:"type:varchar(50)" json:"model"`
	PurchaseDate *time.Time   `json:"purchaseDate,omitempty"`

	SubscriptionID *string    `gorm:"index" json:"subscriptionId,omitempty"`
	BuildingID     *string    `gorm:"index" json:"buildingId,omitempty"`
	AssignedDate   *time.Time `json:"assignedDate,omitempty"`
	CollectionDate *time.Time `json:"collectionDate,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
