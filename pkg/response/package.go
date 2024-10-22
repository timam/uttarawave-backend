package response

import (
	"github.com/timam/uttarawave-backend/models"
	"time"
)

type TVPackageResponse struct {
	ID              string             `json:"id"`
	Type            models.PackageType `json:"type"`
	Name            string             `json:"name"`
	Price           float64            `json:"price"`
	IsActive        bool               `json:"isActive"`
	ConnectionClass string             `json:"connectionClass"`
	ChannelCount    int                `json:"channelCount"`
	TVCount         int                `json:"tvCount"`
	CreatedAt       time.Time          `json:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt"`
}

type InternetPackageResponse struct {
	ID              string             `json:"id"`
	Type            models.PackageType `json:"type"`
	Name            string             `json:"name"`
	Price           float64            `json:"price"`
	IsActive        bool               `json:"isActive"`
	ConnectionClass string             `json:"connectionClass"`
	Bandwidth       int                `json:"bandwidth"`
	BandwidthType   string             `json:"bandwidthType"`
	HasRealIP       bool               `json:"hasRealIP"`
	CreatedAt       time.Time          `json:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt"`
}

func NewTVPackageResponse(pkg *models.Package) TVPackageResponse {
	return TVPackageResponse{
		ID:              pkg.ID,
		Type:            pkg.Type,
		Name:            pkg.Name,
		Price:           pkg.Price,
		IsActive:        pkg.IsActive,
		ConnectionClass: string(pkg.ConnectionClass),
		ChannelCount:    pkg.ChannelCount,
		TVCount:         pkg.TVCount,
		CreatedAt:       pkg.CreatedAt,
		UpdatedAt:       pkg.UpdatedAt,
	}
}

func NewInternetPackageResponse(pkg *models.Package) InternetPackageResponse {
	return InternetPackageResponse{
		ID:              pkg.ID,
		Type:            pkg.Type,
		Name:            pkg.Name,
		Price:           pkg.Price,
		IsActive:        pkg.IsActive,
		ConnectionClass: string(pkg.ConnectionClass),
		Bandwidth:       pkg.Bandwidth,
		BandwidthType:   string(pkg.BandwidthType),
		HasRealIP:       pkg.HasRealIP,
		CreatedAt:       pkg.CreatedAt,
		UpdatedAt:       pkg.UpdatedAt,
	}
}
