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

func NewCustomerHandler() gin.HandlerFunc {
	repo := repositories.NewDynamoDBCustomerRepository()

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

		err := repo.CreateCustomer(&customer)
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

func GetCustomerHandler() gin.HandlerFunc {
	repo := repositories.NewDynamoDBCustomerRepository()

	return func(c *gin.Context) {
		var requestBody struct {
			Mobile string `json:"mobile"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			if err.Error() == "EOF" {
				// No body provided, return all customers
				customers, err := repo.GetAllCustomers()
				if err != nil {
					logger.Error("Failed to get all customers", zap.Error(err))
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get customers"})
					return
				}
				c.JSON(http.StatusOK, customers)
				return
			}
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if requestBody.Mobile == "" {
			// Empty mobile number, return all customers
			customers, err := repo.GetAllCustomers()
			if err != nil {
				logger.Error("Failed to get all customers", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get customers"})
				return
			}
			c.JSON(http.StatusOK, customers)
			return
		}

		// Mobile number provided, return specific customer
		customer, err := repo.GetCustomer(requestBody.Mobile)
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
