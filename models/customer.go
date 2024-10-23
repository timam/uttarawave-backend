package models

import "time"

type CustomerType string

const (
	Individual CustomerType = "Individual"
	Business   CustomerType = "Business"
)

type Customer struct {
	ID                   string       `gorm:"primaryKey" json:"id"`
	Mobile               string       `gorm:"uniqueIndex" json:"mobile"`
	Email                *string      `gorm:"uniqueIndex;null" json:"email,omitempty"`
	Name                 string       `gorm:"type:varchar(100);index" json:"name"`
	Type                 CustomerType `gorm:"type:varchar(20)" json:"type"`
	Address              Address      `gorm:"foreignKey:CustomerID" json:"address"`
	IdentificationNumber *string      `gorm:"type:varchar(50);null;column:identification_number" json:"identificationNumber,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
