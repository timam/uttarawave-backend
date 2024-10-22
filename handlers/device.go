package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/timam/uttarawave-backend/models"
	"github.com/timam/uttarawave-backend/pkg/logger"
	"github.com/timam/uttarawave-backend/pkg/response"
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
			response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
			return
		}

		device.ID = uuid.New().String()

		if device.Brand == "" {
			response.Error(c, http.StatusBadRequest, "Invalid input", "Brand is required")
			return
		}

		now := time.Now()
		if device.PurchaseDate == nil {
			device.PurchaseDate = &now
		}

		err := h.repo.CreateDevice(c.Request.Context(), &device)
		if err != nil {
			logger.Error("Failed to create device", zap.Error(err))
			response.Error(c, http.StatusInternalServerError, "Failed to create device", err.Error())
			return
		}

		logger.Info("Device created successfully", zap.String("id", device.ID))
		deviceResponse := response.NewDeviceResponse(&device)
		response.Success(c, http.StatusCreated, "Device created successfully", deviceResponse)
	}
}

func (h *DeviceHandler) GetDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		device, err := h.repo.GetDeviceByID(c.Request.Context(), id)
		if err != nil {
			logger.Error("Failed to get device", zap.Error(err))
			response.Error(c, http.StatusInternalServerError, "Failed to get device", err.Error())
			return
		}

		deviceResponse := response.NewDeviceResponse(device)
		response.Success(c, http.StatusOK, "Device retrieved successfully", deviceResponse)
	}
}

func (h *DeviceHandler) UpdateDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var updatedDevice models.Device
		if err := c.ShouldBindJSON(&updatedDevice); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
			return
		}

		existingDevice, err := h.repo.GetDeviceByID(c.Request.Context(), id)
		if err != nil {
			logger.Error("Failed to get existing device", zap.Error(err))
			response.Error(c, http.StatusInternalServerError, "Failed to get existing device", err.Error())
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
			response.Error(c, http.StatusInternalServerError, "Failed to update device", err.Error())
			return
		}

		deviceResponse := response.NewDeviceResponse(existingDevice)
		response.Success(c, http.StatusOK, "Device updated successfully", deviceResponse)
	}
}

func (h *DeviceHandler) DeleteDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := h.repo.DeleteDevice(c.Request.Context(), id)
		if err != nil {
			logger.Error("Failed to delete device", zap.Error(err))
			response.Error(c, http.StatusInternalServerError, "Failed to delete device", err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Device deleted successfully", nil)
	}
}

func (h *DeviceHandler) GetAllDevices() gin.HandlerFunc {
	return func(c *gin.Context) {
		devices, err := h.repo.GetAllDevices(c.Request.Context())
		if err != nil {
			logger.Error("Failed to get all devices", zap.Error(err))
			response.Error(c, http.StatusInternalServerError, "Failed to get all devices", err.Error())
			return
		}

		deviceResponses := make([]response.DeviceResponse, len(devices))
		for i, device := range devices {
			deviceResponses[i] = response.NewDeviceResponse(&device)
		}

		response.Success(c, http.StatusOK, "Devices retrieved successfully", deviceResponses)
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
			response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
			return
		}
		if request.AssignmentType == "" || request.AssignmentID == "" {
			response.Error(c, http.StatusBadRequest, "Invalid input", "Assignment type and ID are required")
			return
		}

		err := h.repo.AssignDevice(c.Request.Context(), deviceID, request.AssignmentType, request.AssignmentID)
		if err != nil {
			logger.Error("Failed to assign device", zap.Error(err))
			response.Error(c, http.StatusInternalServerError, "Failed to assign device", err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Device assigned successfully", nil)
	}
}

func (h *DeviceHandler) UnassignDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.Param("id")

		err := h.repo.UnassignDevice(c.Request.Context(), deviceID)
		if err != nil {
			logger.Error("Failed to unassign device", zap.Error(err))
			response.Error(c, http.StatusInternalServerError, "Failed to unassign device", err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Device unassigned successfully", nil)
	}
}

func (h *DeviceHandler) GetDeviceByAssignment() gin.HandlerFunc {
	return func(c *gin.Context) {
		assignmentType := c.Query("assignmentType")
		assignmentID := c.Query("assignmentId")
		if assignmentType == "" || assignmentID == "" {
			response.Error(c, http.StatusBadRequest, "Invalid input", "Assignment type and ID are required")
			return
		}

		device, err := h.repo.GetDeviceByAssignment(c.Request.Context(), assignmentType, assignmentID)
		if err != nil {
			logger.Error("Failed to get device by assignment", zap.Error(err))
			response.Error(c, http.StatusInternalServerError, "Failed to get device", err.Error())
			return
		}

		if device == nil {
			response.Error(c, http.StatusNotFound, "Device not found", "No device found for this assignment")
			return
		}

		deviceResponse := response.NewDeviceResponse(device)
		response.Success(c, http.StatusOK, "Device retrieved successfully", deviceResponse)
	}
}
