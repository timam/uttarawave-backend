package models

import "time"

type PackageType string

const (
	InternetPackage PackageType = "Internet"
	CableTVPackage  PackageType = "CableTV"
)

type BandwidthType string

const (
	Shared    BandwidthType = "shared"
	Dedicated BandwidthType = "dedicated"
)

type Package struct {
	ID       string      `gorm:"primaryKey" json:"id"`
	Type     PackageType `json:"type"`
	Name     string      `gorm:"type:varchar(100)" json:"name"`
	Price    float64     `json:"price"`
	IsActive bool        `json:"isActive"`

	Bandwidth     int           `json:"bandwidth"`
	BandwidthType BandwidthType `json:"bandwidthType"`
	HasRealIP     bool          `json:"hasRealIP"`

	ChannelCount int `json:"channelCount,omitempty"`
	TVCount      int `json:"tvCount,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
