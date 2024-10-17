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
	return &PackageHandler{
		repo: repo,
	}
}

func (h *PackageHandler) CreateInternetPackage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var pkg models.InternetPackage
		if err := c.ShouldBindJSON(&pkg); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		pkg.ID = uuid.New().String()
		err := h.repo.CreateInternetPackage(c.Request.Context(), &pkg)
		if err != nil {
			logger.Error("Failed to create internet package", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create internet package"})
			return
		}

		c.JSON(http.StatusCreated, pkg)
	}
}

func (h *PackageHandler) UpdateInternetPackage() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var pkg models.InternetPackage
		if err := c.ShouldBindJSON(&pkg); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		pkg.ID = id
		err := h.repo.UpdateInternetPackage(c.Request.Context(), &pkg)
		if err != nil {
			logger.Error("Failed to update internet package", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update internet package"})
			return
		}

		c.JSON(http.StatusOK, pkg)
	}
}

func (h *PackageHandler) DeleteInternetPackage() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := h.repo.DeleteInternetPackage(c.Request.Context(), id)
		if err != nil {
			logger.Error("Failed to delete internet package", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete internet package"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Internet package deleted successfully"})
	}
}

func (h *PackageHandler) GetInternetPackage() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		pkg, err := h.repo.GetInternetPackageByID(c.Request.Context(), id)
		if err != nil {
			logger.Error("Failed to get internet package", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get internet package"})
			return
		}

		c.JSON(http.StatusOK, pkg)
	}
}

func (h *PackageHandler) GetAllInternetPackages() gin.HandlerFunc {
	return func(c *gin.Context) {
		packages, err := h.repo.GetAllInternetPackages(c.Request.Context())
		if err != nil {
			logger.Error("Failed to get all internet packages", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all internet packages"})
			return
		}

		c.JSON(http.StatusOK, packages)
	}
}

func (h *PackageHandler) CreateCableTVPackage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var pkg models.CableTVPackage
		if err := c.ShouldBindJSON(&pkg); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		pkg.ID = uuid.New().String()
		err := h.repo.CreateCableTVPackage(c.Request.Context(), &pkg)
		if err != nil {
			logger.Error("Failed to create CableTV package", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create CableTV package"})
			return
		}

		c.JSON(http.StatusCreated, pkg)
	}
}

func (h *PackageHandler) UpdateCableTVPackage() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var pkg models.CableTVPackage
		if err := c.ShouldBindJSON(&pkg); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		pkg.ID = id
		err := h.repo.UpdateCableTVPackage(c.Request.Context(), &pkg)
		if err != nil {
			logger.Error("Failed to update CableTV package", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update CableTV package"})
			return
		}

		c.JSON(http.StatusOK, pkg)
	}
}

func (h *PackageHandler) DeleteCableTVPackage() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := h.repo.DeleteCableTVPackage(c.Request.Context(), id)
		if err != nil {
			logger.Error("Failed to delete CableTV package", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete CableTV package"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "CableTV package deleted successfully"})
	}
}

func (h *PackageHandler) GetCableTVPackage() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		pkg, err := h.repo.GetCableTVPackageByID(c.Request.Context(), id)
		if err != nil {
			logger.Error("Failed to get CableTV package", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get CableTV package"})
			return
		}

		c.JSON(http.StatusOK, pkg)
	}
}

func (h *PackageHandler) GetAllCableTVPackages() gin.HandlerFunc {
	return func(c *gin.Context) {
		packages, err := h.repo.GetAllCableTVPackages(c.Request.Context())
		if err != nil {
			logger.Error("Failed to get all CableTV packages", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all CableTV packages"})
			return
		}

		c.JSON(http.StatusOK, packages)
	}
}
