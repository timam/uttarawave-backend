package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/timam/uttarawave-backend/models"
	"github.com/timam/uttarawave-backend/pkg/logger"
	"github.com/timam/uttarawave-backend/repositories"
	"go.uber.org/zap"
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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		pkg.ID = uuid.New().String()
		err := h.repo.CreatePackage(c.Request.Context(), &pkg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create package"})
			return
		}

		c.JSON(http.StatusCreated, pkg)
	}
}

func (h *PackageHandler) GetPackage() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Package ID is required"})
			return
		}

		pkg, err := h.repo.GetPackageByID(c.Request.Context(), id)
		if err != nil {
			logger.Error("Failed to get package", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get package"})
			return
		}

		if pkg == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Package not found"})
			return
		}

		c.JSON(http.StatusOK, pkg)
	}
}

func (h *PackageHandler) GetPackageByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		pkg, err := h.repo.GetPackageByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Package not found"})
			return
		}

		c.JSON(http.StatusOK, pkg)
	}
}

func (h *PackageHandler) GetAllPackages() gin.HandlerFunc {
	return func(c *gin.Context) {
		packages, err := h.repo.GetAllPackages(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve packages"})
			return
		}

		c.JSON(http.StatusOK, packages)
	}
}

func (h *PackageHandler) UpdatePackage() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var pkg models.Package
		if err := c.ShouldBindJSON(&pkg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		pkg.ID = id
		err := h.repo.UpdatePackage(c.Request.Context(), &pkg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update package"})
			return
		}

		c.JSON(http.StatusOK, pkg)
	}
}

func (h *PackageHandler) DeletePackage() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := h.repo.DeletePackage(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete package"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Package deleted successfully"})
	}
}
