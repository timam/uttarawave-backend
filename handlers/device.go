package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/timam/uttarawave-backend/models"
	"github.com/timam/uttarawave-backend/pkg/logger"
	"github.com/timam/uttarawave-backend/repositories"
	"go.uber.org/zap"
	"net/http"
	"time"
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

		// Validate required fields
		if device.Brand == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Brand is required"})
			return
		}

		// Set date fields only if they are provided
		now := time.Now()
		if device.PurchaseDate == nil {
			device.PurchaseDate = &now
		}
		// AssignedDate and CollectionDate remain nil if not provided

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
		var updatedDevice models.Device
		if err := c.ShouldBindJSON(&updatedDevice); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		existingDevice, err := h.repo.GetDeviceByID(c.Request.Context(), id)
		if err != nil {
			logger.Error("Failed to get existing device", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get existing device"})
			return
		}

		// Update fields
		existingDevice.Brand = updatedDevice.Brand
		existingDevice.Model = updatedDevice.Model
		existingDevice.SerialNumber = updatedDevice.SerialNumber
		existingDevice.Type = updatedDevice.Type
		existingDevice.Usage = updatedDevice.Usage
		existingDevice.Status = updatedDevice.Status

		err = h.repo.UpdateDevice(c.Request.Context(), existingDevice)
		if err != nil {
			logger.Error("Failed to update device", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update device"})
			return
		}

		c.JSON(http.StatusOK, existingDevice)
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

func (h *deviceHandler) AssignDeviceToSubscription() gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.Param("id")
		var request struct {
			SubscriptionID string `json:"subscriptionId"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if request.SubscriptionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Subscription ID is required"})
			return
		}

		err := h.repo.AssignDeviceToSubscription(c.Request.Context(), deviceID, request.SubscriptionID)
		if err != nil {
			logger.Error("Failed to assign device to subscription", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign device to subscription"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Device assigned to subscription successfully"})
	}
}

func (h *deviceHandler) AssignDeviceToBuilding() gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.Param("id")
		var request struct {
			BuildingID string `json:"buildingId"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if request.BuildingID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Building ID is required"})
			return
		}

		err := h.repo.AssignDeviceToBuilding(c.Request.Context(), deviceID, request.BuildingID)
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

func (h *deviceHandler) GetDeviceBySubscriptionID() gin.HandlerFunc {
	return func(c *gin.Context) {
		subscriptionID := c.Query("subscriptionId")
		if subscriptionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Subscription ID is required"})
			return
		}

		device, err := h.repo.GetDeviceBySubscriptionID(c.Request.Context(), subscriptionID)
		if err != nil {
			logger.Error("Failed to get device by subscription ID", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get device"})
			return
		}

		if device == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No device found for this subscription"})
			return
		}

		c.JSON(http.StatusOK, device)
	}
}
