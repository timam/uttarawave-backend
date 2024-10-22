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

type DeviceHandler struct {
	repo repositories.DeviceRepository
}

func NewDeviceHandler(repo repositories.DeviceRepository) *DeviceHandler {
	return &DeviceHandler{
		repo: repo,
	}
}

func (h *DeviceHandler) CreateDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var device models.Device
		if err := c.ShouldBindJSON(&device); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		device.ID = uuid.New().String()

		if device.Brand == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Brand is required"})
			return
		}

		now := time.Now()
		if device.PurchaseDate == nil {
			device.PurchaseDate = &now
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

func (h *DeviceHandler) GetDevice() gin.HandlerFunc {
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

func (h *DeviceHandler) UpdateDevice() gin.HandlerFunc {
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

		existingDevice.Brand = updatedDevice.Brand
		existingDevice.Model = updatedDevice.Model
		existingDevice.SerialNumber = updatedDevice.SerialNumber
		existingDevice.Type = updatedDevice.Type
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

func (h *DeviceHandler) DeleteDevice() gin.HandlerFunc {
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

func (h *DeviceHandler) GetAllDevices() gin.HandlerFunc {
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

func (h *DeviceHandler) AssignDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.Param("id")
		var request struct {
			AssignmentType string `json:"assignmentType"`
			AssignmentID   string `json:"assignmentId"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if request.AssignmentType == "" || request.AssignmentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Assignment type and ID are required"})
			return
		}

		err := h.repo.AssignDevice(c.Request.Context(), deviceID, request.AssignmentType, request.AssignmentID)
		if err != nil {
			logger.Error("Failed to assign device", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign device"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Device assigned successfully"})
	}
}

func (h *DeviceHandler) UnassignDevice() gin.HandlerFunc {
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

func (h *DeviceHandler) GetDeviceByAssignment() gin.HandlerFunc {
	return func(c *gin.Context) {
		assignmentType := c.Query("assignmentType")
		assignmentID := c.Query("assignmentId")
		if assignmentType == "" || assignmentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Assignment type and ID are required"})
			return
		}

		device, err := h.repo.GetDeviceByAssignment(c.Request.Context(), assignmentType, assignmentID)
		if err != nil {
			logger.Error("Failed to get device by assignment", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get device"})
			return
		}

		if device == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No device found for this assignment"})
			return
		}

		c.JSON(http.StatusOK, device)
	}
}
