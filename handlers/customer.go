package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/timam/uttarawave-backend/models"
	"github.com/timam/uttarawave-backend/pkg/logger"
	"github.com/timam/uttarawave-backend/repositories"
	"go.uber.org/zap"
)

type customerHandler struct {
	repo         repositories.CustomerRepository
	buildingRepo repositories.BuildingRepository
}

func NewCustomerHandler(cr repositories.CustomerRepository, br repositories.BuildingRepository) *customerHandler {
	return &customerHandler{
		repo:         cr,
		buildingRepo: br,
	}
}

func generateUniqueID() string {
	return uuid.New().String()
}
func (h *customerHandler) CreateCustomer() gin.HandlerFunc {
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

		// Check if BuildingID is provided
		if customer.BuildingID != "" {
			// Customer is from an existing building
			if customer.Flat == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Flat is required when BuildingID is provided"})
				return
			}
			// Clear individual address fields
			customer.House = ""
			customer.Road = ""
			customer.Block = ""
			customer.Area = ""
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

		err := h.repo.CreateCustomer(c.Request.Context(), &customer)
		if err != nil {
			logger.Error("Failed to save customer data", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save customer data"})
			return
		}

		logger.Info("Customer created successfully",
			zap.String("id", customer.ID),
			zap.String("mobile", customer.Mobile),
		)

		c.JSON(http.StatusCreated, gin.H{"message": "Customer created successfully", "customer": customer})
	}
}

func (h *customerHandler) GetCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		mobile := c.Query("mobile")

		if mobile == "" {
			// No mobile number provided, return all customers
			customers, err := h.repo.GetAllCustomers()
			if err != nil {
				logger.Error("Failed to get all customers", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get customers"})
				return
			}

			// Process each customer to include address information
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
					address := fmt.Sprintf("%s, %s, %s, %s, %s", customer.Flat, building.House, building.Road, building.Block, building.Area)
					customerData["address"] = address
				} else {
					address := fmt.Sprintf("%s, %s, %s, %s", customer.House, customer.Road, customer.Block, customer.Area)
					customerData["address"] = address
				}

				processedCustomers = append(processedCustomers, customerData)
			}

			c.JSON(http.StatusOK, processedCustomers)
			return
		}

		// Mobile number provided, return specific customer
		customer, err := h.repo.GetCustomerByMobile(mobile)
		if err != nil {
			logger.Error("Failed to get customer data", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get customer data"})
			return
		}
		if customer.ID == "" {
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
			address := fmt.Sprintf("%s, %s, %s, %s, %s", customer.Flat, building.House, building.Road, building.Block, building.Area)
			customerData["address"] = address
		} else {
			address := fmt.Sprintf("%s, %s, %s, %s", customer.House, customer.Road, customer.Block, customer.Area)
			customerData["address"] = address
		}

		logger.Info("Customer retrieved successfully", zap.String("mobile", customer.Mobile))
		c.JSON(http.StatusOK, customerData)
	}
}

func (h *customerHandler) UpdateCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var updateData map[string]interface{}

		if err := c.ShouldBindJSON(&updateData); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		mobile, ok := updateData["mobile"].(string)
		if !ok || mobile == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Mobile number must be provided"})
			return
		}

		// Retrieve the existing customer by mobile number.
		existingCustomer, err := h.repo.GetCustomerByMobile(mobile)
		if err != nil {
			logger.Error("Failed to find customer by mobile", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find customer by mobile"})
			return
		}
		if existingCustomer == nil || existingCustomer.ID == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}

		// Update fields based on provided data.
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
		existingCustomer.UpdatedAt = time.Now() // Ensure updated timestamp is set

		err = h.repo.UpdateCustomer(existingCustomer)
		if err != nil {
			logger.Error("Failed to update customer data", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer data"})
			return
		}

		logger.Info("Customer updated successfully",
			zap.String("id", existingCustomer.ID),
			zap.String("mobile", existingCustomer.Mobile),
		)

		c.JSON(http.StatusOK, gin.H{"message": "Customer updated successfully", "customer": existingCustomer})
	}
}
func (h *customerHandler) DeleteCustomer() gin.HandlerFunc {
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
