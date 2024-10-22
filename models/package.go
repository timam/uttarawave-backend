package models

import "time"

type PackageType string

const (
	InternetPackage PackageType = "Internet"
	CableTVPackage  PackageType = "CableTV"
)

type ConnectionClass string

const (
	Home      ConnectionClass = "home"
	Business  ConnectionClass = "business"
	Corporate ConnectionClass = "corporate"
)

type BandwidthType string

const (
	Shared    BandwidthType = "shared"
	Dedicated BandwidthType = "dedicated"
)

type Package struct {
	ID              string          `gorm:"primaryKey" json:"id"`
	Name            string          `gorm:"type:varchar(100)" json:"name"`
	Type            PackageType     `json:"type"`
	Speed           int             `json:"speed,omitempty"` // For Internet packages
	Price           float64         `json:"price"`
	ConnectionType  string          `gorm:"type:varchar(50)" json:"connectionType,omitempty"` // For Internet packages
	ConnectionClass ConnectionClass `json:"connectionClass,omitempty"`                        // For Internet packages
	BandwidthType   BandwidthType   `json:"bandwidthType,omitempty"`                          // For Internet packages
	ChannelLineup   string          `gorm:"type:text" json:"channelLineup,omitempty"`         // For Cable TV packages
	CreatedAt       time.Time       `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time       `gorm:"autoUpdateTime" json:"updatedAt"`
}
