package response

import (
	"github.com/timam/uttarawave-backend/internals/models"
	"time"
)

type PackageResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PackageListResponse struct {
	Items      []interface{}  `json:"items"`
	Pagination PaginationInfo `json:"pagination"`
}

type TVPackageResponse struct {
	ID           string             `json:"id"`
	Type         models.PackageType `json:"type"`
	Name         string             `json:"name"`
	Price        float64            `json:"price"`
	IsActive     bool               `json:"isActive"`
	ChannelCount int                `json:"channelCount"`
	TVCount      int                `json:"tvCount"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`
}

type InternetPackageResponse struct {
	ID            string             `json:"id"`
	Type          models.PackageType `json:"type"`
	Name          string             `json:"name"`
	Price         float64            `json:"price"`
	IsActive      bool               `json:"isActive"`
	Bandwidth     int                `json:"bandwidth"`
	BandwidthType string             `json:"bandwidthType"`
	HasRealIP     bool               `json:"hasRealIP"`
	CreatedAt     time.Time          `json:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt"`
}

func NewPackageResponse(status int, message string, data interface{}) PackageResponse {
	return PackageResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func NewPackageListResponse(packages []interface{}, total int64, page, size int) PackageListResponse {
	return PackageListResponse{
		Items: packages,
		Pagination: PaginationInfo{
			Total: total,
			Page:  page,
			Size:  size,
		},
	}
}
func NewTVPackageResponse(pkg *models.Package) TVPackageResponse {
	return TVPackageResponse{
		ID:           pkg.ID,
		Type:         pkg.Type,
		Name:         pkg.Name,
		Price:        pkg.Price,
		IsActive:     pkg.IsActive,
		ChannelCount: derefInt(pkg.ChannelCount),
		TVCount:      derefInt(pkg.TVCount),
		CreatedAt:    pkg.CreatedAt,
		UpdatedAt:    pkg.UpdatedAt,
	}
}

func NewInternetPackageResponse(pkg *models.Package) InternetPackageResponse {
	return InternetPackageResponse{
		ID:            pkg.ID,
		Type:          pkg.Type,
		Name:          pkg.Name,
		Price:         pkg.Price,
		IsActive:      pkg.IsActive,
		Bandwidth:     derefInt(pkg.Bandwidth),
		BandwidthType: derefBandwidthType(pkg.BandwidthType),
		HasRealIP:     derefBool(pkg.HasRealIP),
		CreatedAt:     pkg.CreatedAt,
		UpdatedAt:     pkg.UpdatedAt,
	}
}

func derefInt(ptr *int) int {
	if ptr != nil {
		return *ptr
	}
	return 0
}

func derefBandwidthType(ptr *models.BandwidthType) string {
	if ptr != nil {
		return string(*ptr)
	}
	return ""
}

func derefBool(ptr *bool) bool {
	if ptr != nil {
		return *ptr
	}
	return false
}
