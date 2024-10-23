package models

import "time"

type Building struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);index" json:"name"`
	Address   Address   `gorm:"foreignKey:BuildingID" json:"address"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
