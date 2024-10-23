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

type BuildingHandler struct {
	repo repositories.BuildingRepository
}

func NewBuildingHandler() *BuildingHandler {
	return &BuildingHandler{
		repo: repositories.NewGormBuildingRepository(),
	}
}
func (h *BuildingHandler) AddBuilding() gin.HandlerFunc {
	return func(c *gin.Context) {
		var building models.Building

		if err := c.ShouldBindJSON(&building); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate required fields
		if building.Name == "" || building.Address.Area == "" || building.Address.Block == "" || building.Address.Road == "" || building.Address.House == "" {
			logger.Warn("Missing required fields")
			c.JSON(http.StatusBadRequest, gin.H{"error": "All fields (Name, Area, Block, Road, House) are required"})
			return
		}

		// Generate a unique ID for the building and address
		building.ID = uuid.New().String()
		building.Address.ID = uuid.New().String()
		building.Address.BuildingID = &building.ID

		err := h.repo.CreateBuilding(c.Request.Context(), &building)
		if err != nil {
			logger.Error("Failed to save building data", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save building data"})
			return
		}

		logger.Info("Building created successfully",
			zap.String("id", building.ID),
			zap.String("name", building.Name),
		)

		c.JSON(http.StatusCreated, gin.H{"message": "Building created successfully", "building": building})
	}
}

func (h *BuildingHandler) DeleteBuilding() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Building ID must be provided"})
			return
		}

		err := h.repo.DeleteBuilding(id)
		if err != nil {
			logger.Error("Failed to delete building", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete building"})
			return
		}

		logger.Info("Building deleted successfully", zap.String("id", id))
		c.JSON(http.StatusOK, gin.H{"message": "Building deleted successfully"})
	}
}

func (h *BuildingHandler) UpdateBuilding() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Building ID must be provided"})
			return
		}

		var updates map[string]interface{}
		if err := c.ShouldBindJSON(&updates); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := h.repo.UpdateBuilding(c.Request.Context(), id, updates)
		if err != nil {
			logger.Error("Failed to update building", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update building"})
			return
		}

		logger.Info("Building updated successfully", zap.String("id", id))
		c.JSON(http.StatusOK, gin.H{"message": "Building updated successfully"})
	}
}

func (h *BuildingHandler) GetBuilding() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Building ID must be provided"})
			return
		}

		building, err := h.repo.GetBuildingByID(c.Request.Context(), id)
		if err != nil {
			logger.Error("Failed to get building", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get building"})
			return
		}

		if building == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Building not found"})
			return
		}

		c.JSON(http.StatusOK, building)
	}
}

func (h *BuildingHandler) GetAllBuildings() gin.HandlerFunc {
	return func(c *gin.Context) {
		buildings, err := h.repo.GetAllBuildings(c.Request.Context())
		if err != nil {
			logger.Error("Failed to get all buildings", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get buildings"})
			return
		}

		c.JSON(http.StatusOK, buildings)
	}
}
