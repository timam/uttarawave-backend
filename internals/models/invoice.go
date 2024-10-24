package models

import "time"

type InvoiceStatus string

const (
	InvoicePending   InvoiceStatus = "PENDING"
	InvoicePaid      InvoiceStatus = "PAID"
	InvoiceOverdue   InvoiceStatus = "OVERDUE"
	InvoiceCancelled InvoiceStatus = "CANCELLED"
)

type Invoice struct {
	ID             string        `gorm:"primaryKey" json:"id"`
	CustomerID     string        `gorm:"index" json:"customerId"`
	SubscriptionID *string       `gorm:"index" json:"subscriptionId,omitempty"`
	Amount         float64       `json:"amount"`
	Status         InvoiceStatus `json:"status"`
	DueDate        time.Time     `json:"dueDate"`
	PaidDate       *time.Time    `json:"paidDate,omitempty"`
	CreatedAt      time.Time     `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time     `gorm:"autoUpdateTime" json:"updatedAt"`
}
