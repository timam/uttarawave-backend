package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/timam/uttarawave-backend/models"
	"github.com/timam/uttarawave-backend/pkg/response"
	"github.com/timam/uttarawave-backend/repositories"
	"net/http"
)

type PackageHandler struct {
	repo repositories.PackageRepository
}

func NewPackageHandler(repo repositories.PackageRepository) *PackageHandler {
	return &PackageHandler{repo: repo}
}

func (h *PackageHandler) CreatePackage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var pkg models.Package
		if err := c.ShouldBindJSON(&pkg); err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
			return
		}

		pkg.ID = uuid.New().String()
		pkg.IsActive = true // Set the package as active by default

		// Validate mandatory fields for all package types
		if pkg.Type == "" || pkg.Name == "" || pkg.Price == 0 || pkg.ConnectionClass == "" {
			response.Error(c, http.StatusBadRequest, "Missing required fields", "type, name, price, and connectionClass are mandatory for package creation")
			return
		}

		// Validate package type and required fields
		switch pkg.Type {
		case models.CableTVPackage:
			if pkg.ChannelCount == 0 || pkg.TVCount == 0 {
				response.Error(c, http.StatusBadRequest, "Missing required fields for Cable TV package", "channelCount and tvCount are required")
				return
			}
		case models.InternetPackage:
			if pkg.Bandwidth == 0 || pkg.BandwidthType == "" {
				response.Error(c, http.StatusBadRequest, "Missing required fields for Internet package", "bandwidth and bandwidthType are required")
				return
			}
		default:
			response.Error(c, http.StatusBadRequest, "Invalid package type", "type must be either CableTV or Internet")
			return
		}

		err := h.repo.CreatePackage(c.Request.Context(), &pkg)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to create package", err.Error())
			return
		}

		var responseData interface{}
		switch pkg.Type {
		case models.CableTVPackage:
			responseData = response.NewTVPackageResponse(&pkg)
		case models.InternetPackage:
			responseData = response.NewInternetPackageResponse(&pkg)
		}

		response.Success(c, http.StatusCreated, "Package created successfully", responseData)
	}
}

func (h *PackageHandler) GetAllPackages() gin.HandlerFunc {
	return func(c *gin.Context) {
		packageType := c.Query("type")

		packages, err := h.repo.GetAllPackages(c.Request.Context(), packageType)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to fetch packages", err.Error())
			return
		}

		if packageType == string(models.CableTVPackage) {
			tvPackages := make([]response.TVPackageResponse, len(packages))
			for i, pkg := range packages {
				tvPackages[i] = response.NewTVPackageResponse(&pkg)
			}
			response.Success(c, http.StatusOK, "TV packages retrieved successfully", tvPackages)
		} else if packageType == string(models.InternetPackage) {
			internetPackages := make([]response.InternetPackageResponse, len(packages))
			for i, pkg := range packages {
				internetPackages[i] = response.NewInternetPackageResponse(&pkg)
			}
			response.Success(c, http.StatusOK, "Internet packages retrieved successfully", internetPackages)
		} else {
			// If no specific type is requested, return all packages
			response.Success(c, http.StatusOK, "Packages retrieved successfully", packages)
		}
	}
}

func (h *PackageHandler) GetPackageByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		pkg, err := h.repo.GetPackageByID(c.Request.Context(), id)
		if err != nil {
			response.Error(c, http.StatusNotFound, "Package not found", err.Error())
			return
		}

		var responseData interface{}
		switch pkg.Type {
		case models.CableTVPackage:
			responseData = response.NewTVPackageResponse(pkg)
		case models.InternetPackage:
			responseData = response.NewInternetPackageResponse(pkg)
		default:
			response.Error(c, http.StatusInternalServerError, "Invalid package type", "Unknown package type")
			return
		}

		response.Success(c, http.StatusOK, "Package retrieved successfully", responseData)
	}
}

func (h *PackageHandler) DeletePackage() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		_, err := h.repo.GetPackageByID(c.Request.Context(), id)
		if err != nil {
			response.Error(c, http.StatusNotFound, "Package not found", err.Error())
			return
		}

		err = h.repo.DeletePackage(c.Request.Context(), id)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to delete package", err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Package deleted successfully", nil)
	}
}
