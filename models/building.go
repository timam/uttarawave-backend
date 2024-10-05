package models

import "time"

type Building struct {
	ID                 string    `gorm:"primaryKey" json:"id"`
	Area               string    `gorm:"type:varchar(100)" json:"area"`
	Block              string    `gorm:"type:varchar(50)" json:"block"`
	Road               string    `gorm:"type:varchar(100)" json:"road"`
	House              string    `gorm:"type:varchar(50)" json:"house"`
	Name               string    `gorm:"type:varchar(100)" json:"name"`
	Devices            []Device  `gorm:"foreignKey:BuildingID" json:"devices,omitempty"`
	HasInternetService bool      `gorm:"default:false" json:"hasInternetService"`
	HasCableTVService  bool      `gorm:"default:false" json:"hasCableTVService"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

type BuildingDevice struct {
	BuildingID string    `gorm:"primaryKey"`
	DeviceID   string    `gorm:"primaryKey"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}
