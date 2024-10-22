package models

import "time"

type PaymentType string

const (
	PaymentIncoming PaymentType = "INCOMING"
	PaymentOutgoing PaymentType = "OUTGOING"
)

type Payment struct {
	ID          string      `gorm:"primaryKey" json:"id"`
	InvoiceID   *string     `gorm:"index" json:"invoiceId,omitempty"`
	CustomerID  *string     `gorm:"index" json:"customerId,omitempty"`
	Amount      float64     `json:"amount"`
	Type        PaymentType `json:"type"`
	Description string      `json:"description"`
	PaidAt      time.Time   `json:"paidAt"`
	CreatedAt   time.Time   `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time   `gorm:"autoUpdateTime" json:"updatedAt"`
}
