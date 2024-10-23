package models

type Address struct {
	ID         string  `gorm:"primaryKey" json:"id"`
	CustomerID *string `gorm:"index" json:"customerId,omitempty"`
	BuildingID *string `gorm:"index" json:"buildingId,omitempty"`

	Flat  string `gorm:"type:varchar(50)" json:"flat,omitempty"`
	House string `gorm:"type:varchar(50)" json:"house,omitempty"`
	Road  string `gorm:"type:varchar(100)" json:"road,omitempty"`
	Block string `gorm:"type:varchar(50)" json:"block,omitempty"`
	Area  string `gorm:"type:varchar(100)" json:"area,omitempty"`
	City  string `gorm:"type:varchar(100)" json:"city,omitempty"`

	Latitude  *float64 `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`
	Longitude *float64 `gorm:"type:decimal(11,8)" json:"longitude,omitempty"`
}
