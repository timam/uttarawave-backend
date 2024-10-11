package models

import "time"

type TransactionType string
type TransactionStatus string

const (
	Cash TransactionType = "CASH"
	// Add more transaction types here in the future
)

const (
	StatusPending   TransactionStatus = "PENDING"
	StatusCompleted TransactionStatus = "COMPLETED"
	StatusFailed    TransactionStatus = "FAILED"
)

type Transaction struct {
	ID             string            `gorm:"primaryKey" json:"id"`
	SubscriptionID string            `gorm:"index" json:"subscriptionId"`
	CustomerID     string            `gorm:"index" json:"customerId"`
	Amount         float64           `json:"amount"`
	Type           TransactionType   `json:"type"`
	Status         TransactionStatus `json:"status"`
	PaidAt         *time.Time        `json:"paidAt,omitempty"`
	CreatedAt      time.Time         `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time         `gorm:"autoUpdateTime" json:"updatedAt"`
}
