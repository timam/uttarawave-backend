package models

type Usage string

const (
	Home     Usage = "home"
	Business Usage = "business"
)

type Type string

const (
	Shared    Type = "shared"
	Dedicated Type = "dedicated"
)

type Package struct {
	Name      string `csv:"Name"`
	Bandwidth string `csv:"Bandwidth"` // e.g., "10Mbps"
	Price     string `csv:"Price"`
	Usage     Usage  `csv:"Usage"`
	Type      Type   `csv:"Type"`
	RealIP    string `csv:"RealIP"`
}
