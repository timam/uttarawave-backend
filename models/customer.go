package models

import "time"

type Customer struct {
	ID         string    `gorm:"primaryKey" json:"id"`
	Mobile     string    `gorm:"uniqueIndex" json:"mobile"`
	Name       string    `gorm:"type:varchar(100)" json:"name"`
	BuildingID string    `gorm:"index;foreignKey:ID" json:"buildingId,omitempty"`
	Flat       string    `gorm:"type:varchar(50)" json:"flat"`
	House      string    `gorm:"type:varchar(50)" json:"house"`
	Road       string    `gorm:"type:varchar(100)" json:"road"`
	Block      string    `gorm:"type:varchar(50)" json:"block"`
	Area       string    `gorm:"type:varchar(100)" json:"area"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
