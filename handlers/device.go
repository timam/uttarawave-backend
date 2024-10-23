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
	"strconv"
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
			c.JSON(http.StatusBadRequest, response.NewDeviceResponse(http.StatusBadRequest, "Invalid input", nil))
			return
		}

		device.ID = uuid.New().String()

		if device.Brand == "" || device.Model == "" || device.SerialNumber == "" || device.Type == "" {
			c.JSON(http.StatusBadRequest, response.NewDeviceResponse(http.StatusBadRequest, "Invalid input", "Brand, Model, SerialNumber, and Type are required"))
			return
		}

		now := time.Now()
		if device.PurchaseDate == nil {
			device.PurchaseDate = &now
		}

		if device.Status == "" {
			device.Status = models.InStock
		}

		if device.Usage == "" {
			device.Usage = models.CompanyUse
		}

		err := h.repo.CreateDevice(c.Request.Context(), &device)
		if err != nil {
			logger.Error("Failed to create device", zap.Error(err))
			c.JSON(http.StatusInternalServerError, response.NewDeviceResponse(http.StatusInternalServerError, "Failed to create device", nil))
			return
		}

		logger.Info("Device created successfully", zap.String("id", device.ID))
		deviceItemResponse := response.NewDeviceItemResponse(&device)
		c.JSON(http.StatusCreated, response.NewDeviceResponse(http.StatusCreated, "Device created successfully", deviceItemResponse))
	}
}

func (h *DeviceHandler) GetDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		device, err := h.repo.GetDeviceByID(c.Request.Context(), id)
		if err != nil {
			logger.Error("Failed to get device", zap.Error(err))
			c.JSON(http.StatusInternalServerError, response.NewDeviceResponse(http.StatusInternalServerError, "Failed to get device", nil))
			return
		}

		deviceItemResponse := response.NewDeviceItemResponse(device)
		c.JSON(http.StatusOK, response.NewDeviceResponse(http.StatusOK, "Device retrieved successfully", deviceItemResponse))
	}
}

func (h *DeviceHandler) UpdateDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var updatedDevice models.Device
		if err := c.ShouldBindJSON(&updatedDevice); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, response.NewDeviceResponse(http.StatusBadRequest, "Invalid input", nil))
			return
		}

		existingDevice, err := h.repo.GetDeviceByID(c.Request.Context(), id)
		if err != nil {
			logger.Error("Failed to get existing device", zap.Error(err))
			c.JSON(http.StatusInternalServerError, response.NewDeviceResponse(http.StatusInternalServerError, "Failed to get existing device", nil))
			return
		}

		existingDevice.Brand = updatedDevice.Brand
		existingDevice.Model = updatedDevice.Model
		existingDevice.SerialNumber = updatedDevice.SerialNumber
		existingDevice.Type = updatedDevice.Type
		existingDevice.Status = updatedDevice.Status
		existingDevice.Usage = updatedDevice.Usage
		existingDevice.PurchasePrice = updatedDevice.PurchasePrice
		existingDevice.PurchaseDate = updatedDevice.PurchaseDate

		err = h.repo.UpdateDevice(c.Request.Context(), existingDevice)
		if err != nil {
			logger.Error("Failed to update device", zap.Error(err))
			c.JSON(http.StatusInternalServerError, response.NewDeviceResponse(http.StatusInternalServerError, "Failed to update device", nil))
			return
		}

		deviceItemResponse := response.NewDeviceItemResponse(existingDevice)
		c.JSON(http.StatusOK, response.NewDeviceResponse(http.StatusOK, "Device updated successfully", deviceItemResponse))
	}
}

func (h *DeviceHandler) DeleteDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := h.repo.DeleteDevice(c.Request.Context(), id)
		if err != nil {
			logger.Error("Failed to delete device", zap.Error(err))
			c.JSON(http.StatusInternalServerError, response.NewDeviceResponse(http.StatusInternalServerError, "Failed to delete device", nil))
			return
		}

		c.JSON(http.StatusOK, response.NewDeviceResponse(http.StatusOK, "Device deleted successfully", nil))
	}
}

func (h *DeviceHandler) GetAllDevices() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

		devices, totalCount, err := h.repo.GetAllDevices(c.Request.Context(), page, pageSize)
		if err != nil {
			logger.Error("Failed to get all devices", zap.Error(err))
			c.JSON(http.StatusInternalServerError, response.NewDeviceResponse(http.StatusInternalServerError, "Failed to get all devices", nil))
			return
		}

		deviceListResponse := response.NewDeviceListResponse(devices, totalCount, page, pageSize)
		c.JSON(http.StatusOK, response.NewDeviceResponse(http.StatusOK, "Devices retrieved successfully", deviceListResponse))
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
			c.JSON(http.StatusBadRequest, response.NewDeviceResponse(http.StatusBadRequest, "Invalid input", nil))
			return
		}
		if request.AssignmentType == "" || request.AssignmentID == "" {
			c.JSON(http.StatusBadRequest, response.NewDeviceResponse(http.StatusBadRequest, "Invalid input", "Assignment type and ID are required"))
			return
		}

		err := h.repo.AssignDevice(c.Request.Context(), deviceID, request.AssignmentType, request.AssignmentID)
		if err != nil {
			logger.Error("Failed to assign device", zap.Error(err))
			c.JSON(http.StatusInternalServerError, response.NewDeviceResponse(http.StatusInternalServerError, "Failed to assign device", nil))
			return
		}

		c.JSON(http.StatusOK, response.NewDeviceResponse(http.StatusOK, "Device assigned successfully", nil))
	}
}

func (h *DeviceHandler) UnassignDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.Param("id")

		err := h.repo.UnassignDevice(c.Request.Context(), deviceID)
		if err != nil {
			logger.Error("Failed to unassign device", zap.Error(err))
			c.JSON(http.StatusInternalServerError, response.NewDeviceResponse(http.StatusInternalServerError, "Failed to unassign device", nil))
			return
		}

		c.JSON(http.StatusOK, response.NewDeviceResponse(http.StatusOK, "Device unassigned successfully", nil))
	}
}

func (h *DeviceHandler) GetDeviceByAssignment() gin.HandlerFunc {
	return func(c *gin.Context) {
		assignmentType := c.Query("assignmentType")
		assignmentID := c.Query("assignmentId")
		if assignmentType == "" || assignmentID == "" {
			c.JSON(http.StatusBadRequest, response.NewDeviceResponse(http.StatusBadRequest, "Invalid input", "Assignment type and ID are required"))
			return
		}

		device, err := h.repo.GetDeviceByAssignment(c.Request.Context(), assignmentType, assignmentID)
		if err != nil {
			logger.Error("Failed to get device by assignment", zap.Error(err))
			c.JSON(http.StatusInternalServerError, response.NewDeviceResponse(http.StatusInternalServerError, "Failed to get device", nil))
			return
		}

		if device == nil {
			c.JSON(http.StatusNotFound, response.NewDeviceResponse(http.StatusNotFound, "Device not found", "No device found for this assignment"))
			return
		}

		deviceItemResponse := response.NewDeviceItemResponse(device)
		c.JSON(http.StatusOK, response.NewDeviceResponse(http.StatusOK, "Device retrieved successfully", deviceItemResponse))
	}
}
