package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/timam/uttarawave-backend/models"
	"github.com/timam/uttarawave-backend/pkg/logger"
	"github.com/timam/uttarawave-backend/repositories"
	"go.uber.org/zap"
)

type CustomerHandler struct {
	repo             repositories.CustomerRepository
	buildingRepo     repositories.BuildingRepository
	subscriptionRepo repositories.SubscriptionRepository
	deviceRepo       repositories.DeviceRepository
}

func NewCustomerHandler(
	cr repositories.CustomerRepository,
	br repositories.BuildingRepository,
	sr repositories.SubscriptionRepository,
	dr repositories.DeviceRepository) *CustomerHandler {
	return &CustomerHandler{
		repo:             cr,
		buildingRepo:     br,
		subscriptionRepo: sr,
		deviceRepo:       dr,
	}
}

func generateUniqueID() string {
	return uuid.New().String()
}

func (h *CustomerHandler) CreateCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var customer models.Customer

		if err := c.ShouldBindJSON(&customer); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate required fields
		if customer.Mobile == "" || customer.Name == "" {
			logger.Warn("Missing required fields")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Mobile and Name are required fields"})
			return
		}

		var buildingDetails *models.Building
		var err error

		// Check if BuildingID is provided
		if customer.BuildingID != "" {
			// Customer is from an existing building
			if customer.Flat == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Flat is required when BuildingID is provided"})
				return
			}
			// Get building details
			buildingDetails, err = h.buildingRepo.GetBuildingDetails(c.Request.Context(), customer.BuildingID)
			if err != nil {
				logger.Error("Failed to get building details", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get building details"})
				return
			}
			// Set address fields from building details
			customer.House = buildingDetails.House
			customer.Road = buildingDetails.Road
			customer.Block = buildingDetails.Block
			customer.Area = buildingDetails.Area
		} else {
			// Customer is not from an existing building
			// Validate address fields
			if customer.House == "" || customer.Road == "" || customer.Block == "" || customer.Area == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "All address fields (House, Road, Block, Area) are required when BuildingID is not provided"})
				return
			}
		}

		// Generate a unique ID for the customer
		customer.ID = generateUniqueID()

		err = h.repo.CreateCustomer(c.Request.Context(), &customer)
		if err != nil {
			logger.Error("Failed to save customer data", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save customer data"})
			return
		}

		logger.Info("Customer created successfully",
			zap.String("id", customer.ID),
			zap.String("mobile", customer.Mobile),
		)

		response := gin.H{
			"message":  "Customer created successfully",
			"customer": customer,
		}

		c.JSON(http.StatusCreated, response)
	}
}

func (h *CustomerHandler) GetCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		mobile := c.Query("mobile")

		if mobile == "" {
			h.GetAllCustomers()(c)
			return
		}

		customer, err := h.repo.GetCustomerByMobile(mobile)
		if err != nil {
			logger.Error("Failed to get customer data", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get customer data"})
			return
		}
		if customer == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}

		customerData := gin.H{
			"id":     customer.ID,
			"mobile": customer.Mobile,
			"name":   customer.Name,
		}

		if customer.BuildingID != "" {
			building, err := h.buildingRepo.GetBuildingByID(c.Request.Context(), customer.BuildingID)
			if err != nil {
				logger.Error("Failed to get building data", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get building data"})
				return
			}
			customerData["address"] = gin.H{
				"flat":  customer.Flat,
				"house": building.House,
				"road":  building.Road,
				"block": building.Block,
				"area":  building.Area,
			}
		} else {
			customerData["address"] = gin.H{
				"flat":  customer.Flat,
				"house": customer.House,
				"road":  customer.Road,
				"block": customer.Block,
				"area":  customer.Area,
			}
		}

		logger.Info("Customer retrieved successfully", zap.String("mobile", customer.Mobile))
		c.JSON(http.StatusOK, customerData)
	}
}

func (h *CustomerHandler) GetAllCustomers() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))

		customers, totalCount, err := h.repo.GetCustomersPaginated(page, pageSize)
		if err != nil {
			logger.Error("Failed to get customers", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get customers"})
			return
		}

		var processedCustomers []gin.H
		for _, customer := range customers {
			customerData := gin.H{
				"id":     customer.ID,
				"mobile": customer.Mobile,
				"name":   customer.Name,
			}

			if customer.BuildingID != "" {
				building, err := h.buildingRepo.GetBuildingByID(c.Request.Context(), customer.BuildingID)
				if err != nil {
					logger.Error("Failed to get building data", zap.Error(err))
					continue
				}
				customerData["address"] = gin.H{
					"flat":  customer.Flat,
					"house": building.House,
					"road":  building.Road,
					"block": building.Block,
					"area":  building.Area,
				}
			} else {
				customerData["address"] = gin.H{
					"flat":  customer.Flat,
					"house": customer.House,
					"road":  customer.Road,
					"block": customer.Block,
					"area":  customer.Area,
				}
			}

			processedCustomers = append(processedCustomers, customerData)
		}

		response := gin.H{
			"customers":  processedCustomers,
			"totalCount": totalCount,
			"page":       page,
			"pageSize":   pageSize,
		}

		logger.Info("Retrieved customers", zap.Int("count", len(processedCustomers)), zap.Int("page", page), zap.Int("pageSize", pageSize))
		c.JSON(http.StatusOK, response)
	}
}

func (h *CustomerHandler) UpdateCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		logger.Info("UpdateCustomer called", zap.String("id", id))

		if id == "" {
			logger.Warn("Customer ID not provided")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Customer ID must be provided"})
			return
		}

		var updateData map[string]interface{}
		if err := c.ShouldBindJSON(&updateData); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		logger.Info("Update data received", zap.Any("updateData", updateData))

		// Retrieve the existing customer by ID
		existingCustomer, err := h.repo.GetCustomer(id)
		if err != nil {
			logger.Error("Failed to find customer by ID", zap.Error(err), zap.String("id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find customer"})
			return
		}
		if existingCustomer == nil {
			logger.Warn("Customer not found", zap.String("id", id))
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}
		logger.Info("Existing customer found", zap.String("id", existingCustomer.ID), zap.String("name", existingCustomer.Name))

		// Prevent updating the mobile number
		if _, exists := updateData["mobile"]; exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Mobile number cannot be updated"})
			return
		}

		// Update fields based on provided data
		if name, ok := updateData["name"].(string); ok {
			existingCustomer.Name = name
		}
		if buildingId, ok := updateData["buildingId"].(string); ok {
			existingCustomer.BuildingID = buildingId
		}
		if flat, ok := updateData["flat"].(string); ok {
			existingCustomer.Flat = flat
		}
		if house, ok := updateData["house"].(string); ok {
			existingCustomer.House = house
		}
		if road, ok := updateData["road"].(string); ok {
			existingCustomer.Road = road
		}
		if block, ok := updateData["block"].(string); ok {
			existingCustomer.Block = block
		}
		if area, ok := updateData["area"].(string); ok {
			existingCustomer.Area = area
		}
		existingCustomer.UpdatedAt = time.Now()

		err = h.repo.UpdateCustomer(existingCustomer)
		if err != nil {
			logger.Error("Failed to update customer data", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer data"})
			return
		}

		logger.Info("Customer updated successfully", zap.String("id", existingCustomer.ID), zap.String("name", existingCustomer.Name))
		c.JSON(http.StatusOK, gin.H{"message": "Customer updated successfully", "customer": existingCustomer})
	}
}

func (h *CustomerHandler) DeleteCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		mobile := c.Query("mobile")

		if mobile == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Mobile number must be provided"})
			return
		}

		customer, err := h.repo.GetCustomerByMobile(mobile)
		if err != nil {
			logger.Error("Failed to find customer by mobile", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find customer by mobile"})
			return
		}
		if customer.ID == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}

		if err := h.repo.DeleteCustomer(customer.ID); err != nil {
			logger.Error("Failed to delete customer by mobile", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
	}
}

//func (h *CustomerHandler) GetAllCustomersFullDetails() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
//		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))
//
//		customers, totalCount, err := h.repo.GetCustomersPaginated(page, pageSize)
//		if err != nil {
//			logger.Error("Failed to get customers", zap.Error(err))
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get customers"})
//			return
//		}
//
//		var fullCustomersDetails []gin.H
//		for _, customer := range customers {
//			subscriptions, err := h.subscriptionRepo.GetSubscriptionsByCustomerID(c.Request.Context(), customer.ID)
//			if err != nil {
//				logger.Error("Failed to get subscriptions", zap.Error(err), zap.String("customerID", customer.ID))
//				continue
//			}
//
//			var fullSubscriptions []gin.H
//			for _, sub := range subscriptions {
//				device, err := h.deviceRepo.GetDeviceBySubscriptionID(c.Request.Context(), sub.ID)
//				if err != nil {
//					logger.Error("Failed to get device", zap.Error(err), zap.String("subscriptionID", sub.ID))
//					continue
//				}
//
//				fullSub := gin.H{
//					"subscription": sub,
//					"device":       device,
//				}
//				fullSubscriptions = append(fullSubscriptions, fullSub)
//			}
//
//			var buildingDetails gin.H
//			if customer.BuildingID != "" {
//				building, err := h.buildingRepo.GetBuildingByID(c.Request.Context(), customer.BuildingID)
//				if err != nil {
//					logger.Error("Failed to get building details", zap.Error(err), zap.String("buildingID", customer.BuildingID))
//				} else {
//					buildingDetails = gin.H{
//						"id":    building.ID,
//						"name":  building.Name,
//						"house": building.House,
//						"road":  building.Road,
//						"block": building.Block,
//						"area":  building.Area,
//					}
//				}
//			}
//
//			customerDetails := gin.H{
//				"customer":      customer,
//				"subscriptions": fullSubscriptions,
//				"building":      buildingDetails,
//			}
//			fullCustomersDetails = append(fullCustomersDetails, customerDetails)
//		}
//
//		response := gin.H{
//			"customers":  fullCustomersDetails,
//			"totalCount": totalCount,
//			"page":       page,
//			"pageSize":   pageSize,
//		}
//
//		c.JSON(http.StatusOK, response)
//	}
//}
