package models

import "time"

type IncomeType string

const (
	SubscriptionPayment IncomeType = "SUBSCRIPTION_PAYMENT"
	OtherIncome         IncomeType = "OTHER_INCOME"
)

type Income struct {
	ID             string     `gorm:"primaryKey" json:"id"`
	SubscriptionID *string    `gorm:"index" json:"subscriptionId,omitempty"`
	CustomerID     string     `gorm:"index" json:"customerId"`
	Amount         float64    `json:"amount"`
	Type           IncomeType `json:"type"`
	Description    string     `json:"description"`
	ReceivedAt     time.Time  `json:"receivedAt"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
}
