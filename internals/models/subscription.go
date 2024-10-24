package models

import "time"

type SubscriptionType string

const (
	Internet SubscriptionType = "Internet"
	CableTV  SubscriptionType = "CableTV"
)

type Subscription struct {
	ID              string    `gorm:"primaryKey" json:"id"`
	CustomerID      string    `gorm:"index" json:"customerId"`
	PackageID       string    `gorm:"index" json:"packageId"`
	PackagePrice    float64   `json:"packagePrice"`
	MonthlyDiscount float64   `json:"monthlyDiscount"`
	Status          string    `json:"status"`
	StartDate       time.Time `json:"startDate"`
	RenewalDate     time.Time `json:"renewalDate"`
	PaidUntil       time.Time `json:"paidUntil"`
	DueAmount       string    `json:"dueAmount"`
	DeviceID        string    `gorm:"index" json:"deviceId,omitempty"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
