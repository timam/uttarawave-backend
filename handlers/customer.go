package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/timam/uttarawave-finance-backend/models"
	"github.com/timam/uttarawave-finance-backend/repositories"
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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate required fields
		if customer.Mobile == "" || customer.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Mobile and Name are required fields"})
			return
		}

		err := repo.CreateCustomer(&customer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save customer data"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Customer created successfully", "customer": customer})
	}
}
