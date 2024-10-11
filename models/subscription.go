package models

import "time"

type SubscriptionType string

const (
	Internet SubscriptionType = "Internet"
	CableTV  SubscriptionType = "CableTV"
)

type Subscription struct {
	ID           string           `gorm:"primaryKey" json:"id"`
	CustomerID   string           `gorm:"index" json:"customerId"`
	Type         SubscriptionType `json:"type"`
	PackageID    string           `gorm:"index" json:"packageId"`
	PackageName  string           `json:"packageName"`
	PackagePrice string           `json:"packagePrice"`
	Discount     string           `json:"discount"`
	Status       string           `json:"status"`
	StartDate    time.Time        `json:"startDate"`
	RenewalDate  time.Time        `json:"renewalDate"`
	PaidUntil    time.Time        `json:"paidUntil"`
	DeviceID     string           `gorm:"index" json:"deviceId,omitempty"`
	CreatedAt    time.Time        `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time        `gorm:"autoUpdateTime" json:"updatedAt"`
}
