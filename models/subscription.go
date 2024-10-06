package models

import "time"

type SubscriptionType string

const (
	Internet SubscriptionType = "Internet"
	CableTV  SubscriptionType = "CableTV"
)

type Subscription struct {
	ID             string           `gorm:"primaryKey" json:"id"`
	CustomerID     string           `gorm:"index" json:"customerId"`
	Type           SubscriptionType `json:"type"`
	PackageName    string           `json:"packageName"`
	PackagePrice   int              `json:"packagePrice"`
	PackageVersion string           `json:"packageVersion"`
	Discount       int              `json:"discount"`
	Status         string           `json:"status"`
	StartDate      time.Time        `json:"startDate"`
	RenewalDate    time.Time        `json:"renewalDate"`
	DeviceID       string           `gorm:"index" json:"deviceId,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
