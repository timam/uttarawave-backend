package models

type PackageType string

const (
	Home     PackageType = "home"
	Business PackageType = "business"
)

type ConnectionType string

const (
	Shared    ConnectionType = "shared"
	Dedicated ConnectionType = "dedicated"
)

type Package struct {
	Name       string         `csv:"Name"`
	Speed      string         `csv:"Speed"` // e.g., "10Mbps"
	Price      string         `csv:"Price"`
	Connection PackageType    `csv:"Connection"`
	Type       ConnectionType `csv:"Type"`
	RealIP     string         `csv:"RealIP"`
}
