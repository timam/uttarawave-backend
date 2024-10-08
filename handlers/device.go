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

type deviceHandler struct {
	repo         repositories.DeviceRepository
	buildingRepo repositories.BuildingRepository
}

func NewDeviceHandler(repo repositories.DeviceRepository, buildingRepo repositories.BuildingRepository) *deviceHandler {
	return &deviceHandler{
		repo:         repo,
		buildingRepo: buildingRepo,
	}
}

func (h *deviceHandler) CreateDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var device models.Device
		if err := c.ShouldBindJSON(&device); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Generate a unique ID for the device
		device.ID = uuid.New().String()

		// If BuildingID is provided, verify that the building exists
		if device.BuildingID != nil && *device.BuildingID != "" {
			building, err := h.buildingRepo.GetBuildingByID(c.Request.Context(), *device.BuildingID)
			if err != nil || building == nil {
				logger.Error("Invalid BuildingID provided", zap.Error(err))
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid BuildingID provided"})
				return
			}
		}

		err := h.repo.CreateDevice(c.Request.Context(), &device)
		if err != nil {
			logger.Error("Failed to create device", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create device"})
			return
		}

		logger.Info("Device created successfully", zap.String("id", device.ID))
		c.JSON(http.StatusCreated, gin.H{"message": "Device created successfully", "device": device})
	}
}

func (h *deviceHandler) GetDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		device, err := h.repo.GetDeviceByID(c.Request.Context(), id)
		if err != nil {
			logger.Error("Failed to get device", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get device"})
			return
		}

		c.JSON(http.StatusOK, device)
	}
}

func (h *deviceHandler) UpdateDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var device models.Device
		if err := c.ShouldBindJSON(&device); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		device.ID = id
		err := h.repo.UpdateDevice(c.Request.Context(), &device)
		if err != nil {
			logger.Error("Failed to update device", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update device"})
			return
		}

		c.JSON(http.StatusOK, device)
	}
}

func (h *deviceHandler) DeleteDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := h.repo.DeleteDevice(c.Request.Context(), id)
		if err != nil {
			logger.Error("Failed to delete device", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete device"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Device deleted successfully"})
	}
}

func (h *deviceHandler) GetAllDevices() gin.HandlerFunc {
	return func(c *gin.Context) {
		devices, err := h.repo.GetAllDevices(c.Request.Context())
		if err != nil {
			logger.Error("Failed to get all devices", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all devices"})
			return
		}

		c.JSON(http.StatusOK, devices)
	}
}

func (h *deviceHandler) AssignDeviceToCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.Param("id")
		customerID := c.Query("customerId")
		subscriptionID := c.Query("subscriptionId")
		if customerID == "" || subscriptionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Customer ID and Subscription ID are required"})
			return
		}

		err := h.repo.AssignDeviceToCustomer(c.Request.Context(), deviceID, customerID, subscriptionID)
		if err != nil {
			logger.Error("Failed to assign device to customer", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign device to customer"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Device assigned to customer successfully"})
	}
}

func (h *deviceHandler) AssignDeviceToBuilding() gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.Param("id")
		buildingID := c.Query("buildingId")
		if buildingID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Building ID is required"})
			return
		}

		err := h.repo.AssignDeviceToBuilding(c.Request.Context(), deviceID, buildingID)
		if err != nil {
			logger.Error("Failed to assign device to building", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign device to building"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Device assigned to building successfully"})
	}
}

func (h *deviceHandler) UnassignDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.Param("id")

		err := h.repo.UnassignDevice(c.Request.Context(), deviceID)
		if err != nil {
			logger.Error("Failed to unassign device", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unassign device"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Device unassigned successfully"})
	}
}
