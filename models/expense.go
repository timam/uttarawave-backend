package models

import "time"

type ExpenseType string

const (
	OperationalExpense ExpenseType = "OPERATIONAL_EXPENSE"
	CapitalExpense     ExpenseType = "CAPITAL_EXPENSE"
)

type Expense struct {
	ID          string      `gorm:"primaryKey" json:"id"`
	Amount      float64     `json:"amount"`
	Type        ExpenseType `json:"type"`
	Description string      `json:"description"`
	PaidAt      time.Time   `json:"paidAt"`
	CreatedAt   time.Time   `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time   `gorm:"autoUpdateTime" json:"updatedAt"`
}
