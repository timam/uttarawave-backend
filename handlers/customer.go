package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/timam/uttarawave-backend/models"
	"github.com/timam/uttarawave-backend/pkg/logger"
	"github.com/timam/uttarawave-backend/repositories"
	"go.uber.org/zap"
	"net/http"
)

type CustomerHandler struct {
	repo repositories.CustomerRepository
}

func NewCustomerHandler() *CustomerHandler {
	return &CustomerHandler{
		repo: repositories.NewGormCustomerRepository(),
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

		err := h.repo.CreateCustomer(&customer)
		if err != nil {
			logger.Error("Failed to save customer data", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save customer data"})
			return
		}

		logger.Info("Customer created successfully",
			zap.String("mobile", customer.Mobile),
		)

		c.JSON(http.StatusCreated, gin.H{"message": "Customer created successfully", "customer": customer})
	}
}

func (h *CustomerHandler) GetCustomer() gin.HandlerFunc {
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
			c.JSON(http.StatusOK, customers)
			return
		}

		// Mobile number provided, return specific customer
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

		logger.Info("Customer retrieved successfully",
			zap.String("mobile", customer.Mobile),
		)

		c.JSON(http.StatusOK, customer)
	}
}

func (h *CustomerHandler) UpdateCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var customer models.Customer

		if err := c.ShouldBindJSON(&customer); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := h.repo.UpdateCustomer(&customer)
		if err != nil {
			logger.Error("Failed to update customer data", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer data"})
			return
		}

		logger.Info("Customer updated successfully",
			zap.String("mobile", customer.Mobile),
		)

		c.JSON(http.StatusOK, gin.H{"message": "Customer updated successfully", "customer": customer})
	}
}

func (h *CustomerHandler) DeleteCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := h.repo.DeleteCustomer(id)
		if err != nil {
			logger.Error("Failed to delete customer", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer"})
			return
		}

		logger.Info("Customer deleted successfully",
			zap.String("id", id),
		)

		c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
	}
}
