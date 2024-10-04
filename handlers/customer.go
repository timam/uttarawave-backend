package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/timam/uttarawave-finance-backend/models"
	"github.com/timam/uttarawave-finance-backend/pkg/logger"
	"github.com/timam/uttarawave-finance-backend/repositories"
	"go.opentelemetry.io/otel/trace" // Add this line
	"go.uber.org/zap"
	"net/http"
)

type CustomerHandler struct {
	repo repositories.CustomerRepository
}

func NewCustomerHandler() gin.HandlerFunc {
	repo := repositories.NewDynamoDBCustomerRepository()

	return func(c *gin.Context) {
		span, exists := c.Get("span")
		if !exists {
			logger.Error("Span not found in context")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

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
			zap.String("traceID", span.(trace.Span).SpanContext().TraceID().String()),
		)

		c.JSON(http.StatusCreated, gin.H{"message": "Customer created successfully", "customer": customer})
	}
}
