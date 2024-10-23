package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/timam/uttarawave-backend/models"
	"github.com/timam/uttarawave-backend/pkg/logger"
	"github.com/timam/uttarawave-backend/repositories"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type CustomerHandler struct {
	repo         repositories.CustomerRepository
	buildingRepo repositories.BuildingRepository
}

func NewCustomerHandler(
	cr repositories.CustomerRepository,
	br repositories.BuildingRepository) *CustomerHandler {
	return &CustomerHandler{
		repo:         cr,
		buildingRepo: br,
	}
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

		// Generate a unique ID for the customer
		customer.ID = uuid.New().String()

		if customer.Address.BuildingID != nil {
			building, err := h.buildingRepo.GetBuildingByID(c.Request.Context(), *customer.Address.BuildingID)
			if err != nil {
				logger.Error("Failed to get building details", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get building details"})
				return
			}
			// Use the building's address
			customer.Address = building.Address
			customer.Address.ID = uuid.New().String()
			customer.Address.CustomerID = &customer.ID
		} else {
			// Validate address fields for individual address
			if customer.Address.House == "" || customer.Address.Road == "" || customer.Address.Block == "" || customer.Address.Area == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "All address fields (House, Road, Block, Area) are required when BuildingID is not provided"})
				return
			}
			customer.Address.ID = uuid.New().String()
			customer.Address.CustomerID = &customer.ID
		}

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

		c.JSON(http.StatusOK, customer)
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

		response := gin.H{
			"customers":  customers,
			"totalCount": totalCount,
			"page":       page,
			"pageSize":   pageSize,
		}

		logger.Info("Retrieved customers", zap.Int("count", len(customers)), zap.Int("page", page), zap.Int("pageSize", pageSize))
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

		var updateData models.Customer
		if err := c.ShouldBindJSON(&updateData); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		existingCustomer, err := h.repo.GetCustomer(id)
		if err != nil {
			logger.Error("Failed to find customer", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find customer"})
			return
		}
		if existingCustomer == nil {
			logger.Warn("Customer not found", zap.String("id", id))
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}

		// Prevent updating the mobile number
		if updateData.Mobile != "" && updateData.Mobile != existingCustomer.Mobile {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Mobile number cannot be updated"})
			return
		}

		// Update fields
		existingCustomer.Name = updateData.Name
		existingCustomer.Email = updateData.Email
		existingCustomer.Type = updateData.Type
		existingCustomer.IdentificationNumber = updateData.IdentificationNumber

		// Handle address update
		if updateData.Address.BuildingID != nil {
			building, err := h.buildingRepo.GetBuildingByID(c.Request.Context(), *updateData.Address.BuildingID)
			if err != nil {
				logger.Error("Failed to get building details", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get building details"})
				return
			}
			existingCustomer.Address = building.Address
			existingCustomer.Address.ID = uuid.New().String()
			existingCustomer.Address.CustomerID = &existingCustomer.ID
			existingCustomer.Address.BuildingID = updateData.Address.BuildingID
		} else {
			// Validate address fields for individual address
			if updateData.Address.House == "" || updateData.Address.Road == "" || updateData.Address.Block == "" || updateData.Address.Area == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "All address fields (House, Road, Block, Area) are required when BuildingID is not provided"})
				return
			}
			existingCustomer.Address = updateData.Address
			existingCustomer.Address.ID = uuid.New().String()
			existingCustomer.Address.CustomerID = &existingCustomer.ID
			existingCustomer.Address.BuildingID = nil
		}

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
