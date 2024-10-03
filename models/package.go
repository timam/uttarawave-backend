package models

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

type InternetPackage struct {
	PackageName     string          `csv:"PackageName"`
	Bandwidth       string          `csv:"Bandwidth"` // e.g., "10Mbps"
	Price           string          `csv:"Price"`
	ConnectionClass ConnectionClass `csv:"ConnectionClass"`
	BandwidthType   BandwidthType   `csv:"BandwidthType"`
	RealIP          string          `csv:"RealIP"`
}

type CableTVPackage struct {
	PackageName string `csv:"PackageName"`
	Price       string `csv:"Price"`
	TVCount     string `csv:"TVCount"`
}
