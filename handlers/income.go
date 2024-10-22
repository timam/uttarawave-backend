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

type IncomeHandler struct {
	incomeRepo       repositories.IncomeRepository
	subscriptionRepo repositories.SubscriptionRepository
}

func NewIncomeHandler(ir repositories.IncomeRepository, sr repositories.SubscriptionRepository) *IncomeHandler {
	return &IncomeHandler{
		incomeRepo:       ir,
		subscriptionRepo: sr,
	}
}

func (h *IncomeHandler) CreateIncome() gin.HandlerFunc {
	return func(c *gin.Context) {
		var income models.Income
		if err := c.ShouldBindJSON(&income); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		income.ID = uuid.New().String()
		income.ReceivedAt = time.Now()

		if income.Type == models.SubscriptionPayment && income.SubscriptionID != nil {
			subscription, err := h.subscriptionRepo.GetSubscription(c.Request.Context(), *income.SubscriptionID)
			if err != nil {
				logger.Error("Failed to get subscription", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subscription"})
				return
			}

			// Update subscription (similar to existing ProcessCashTransaction logic)
			// ...

			err = h.subscriptionRepo.UpdateSubscription(c.Request.Context(), subscription)
			if err != nil {
				logger.Error("Failed to update subscription", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription"})
				return
			}
		}

		err := h.incomeRepo.CreateIncome(c.Request.Context(), &income)
		if err != nil {
			logger.Error("Failed to create income", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create income"})
			return
		}

		c.JSON(http.StatusCreated, income)
	}
}

// Add other income handler methods here (Get, GetAll, etc.)
