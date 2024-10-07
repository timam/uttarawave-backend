package models

import (
	"time"
)

type ConnectionClass string
type BandwidthType string

const (
	Home      ConnectionClass = "home"
	Business  ConnectionClass = "business"
	Corporate ConnectionClass = "corporate"
)

const (
	Shared    BandwidthType = "shared"
	Dedicated BandwidthType = "dedicated"
)

type InternetPackage struct {
	ID              string          `gorm:"primaryKey" json:"id"`
	PackageName     string          `gorm:"uniqueIndex" json:"packageName"`
	Bandwidth       string          `json:"bandwidth"`
	Price           string          `json:"price"`
	ConnectionClass ConnectionClass `json:"connectionClass"`
	BandwidthType   BandwidthType   `json:"bandwidthType"`
	RealIP          string          `json:"realIP"`
	IsActive        bool            `gorm:"default:true" json:"isActive"`
	CreatedAt       time.Time       `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time       `gorm:"autoUpdateTime" json:"updatedAt"`
}

type CableTVPackage struct {
	ID              string          `gorm:"primaryKey" json:"id"`
	PackageName     string          `gorm:"uniqueIndex" json:"packageName"`
	Price           string          `json:"price"`
	ConnectionClass ConnectionClass `json:"connectionClass"`
	TVCount         string          `json:"tvCount"`
	IsActive        bool            `gorm:"default:true" json:"isActive"`
	CreatedAt       time.Time       `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time       `gorm:"autoUpdateTime" json:"updatedAt"`
}
