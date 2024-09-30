package models

type ConnectionType string

const (
	Home      ConnectionType = "home"
	Business  ConnectionType = "business"
	Corporate ConnectionType = "corporate"
)

type BandwidthType string

const (
	Shared    BandwidthType = "shared"
	Dedicated BandwidthType = "dedicated"
)

type InternetPackage struct {
	PackageName    string         `csv:"PackageName"`
	Bandwidth      string         `csv:"Bandwidth"` // e.g., "10Mbps"
	Price          string         `csv:"Price"`
	ConnectionType ConnectionType `csv:"ConnectionType"`
	BandwidthType  BandwidthType  `csv:"BandwidthType"`
	RealIP         string         `csv:"RealIP"`
}

type CableTVPackage struct {
	Name    string `csv:"PackageName"`
	Price   string `csv:"Price"`
	TVCount string `csv:"TVCount"`
}
