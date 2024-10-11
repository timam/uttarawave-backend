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
	"time"
)

type transactionHandler struct {
	transactionRepo  repositories.TransactionRepository
	subscriptionRepo repositories.SubscriptionRepository
}

func NewTransactionHandler(tr repositories.TransactionRepository, sr repositories.SubscriptionRepository) *transactionHandler {
	return &transactionHandler{
		transactionRepo:  tr,
		subscriptionRepo: sr,
	}
}

func (h *transactionHandler) ProcessCashTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var transactionInfo struct {
			SubscriptionID string  `json:"subscriptionId"`
			Amount         float64 `json:"amount"`
		}

		if err := c.ShouldBindJSON(&transactionInfo); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		subscription, err := h.subscriptionRepo.GetSubscription(c.Request.Context(), transactionInfo.SubscriptionID)
		if err != nil {
			logger.Error("Failed to get subscription", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subscription"})
			return
		}

		now := time.Now()
		transaction := &models.Transaction{
			ID:             uuid.New().String(),
			SubscriptionID: transactionInfo.SubscriptionID,
			CustomerID:     subscription.CustomerID,
			Amount:         transactionInfo.Amount,
			Type:           models.Cash,
			Status:         models.StatusCompleted,
			PaidAt:         &now,
		}

		err = h.transactionRepo.CreateTransaction(c.Request.Context(), transaction)
		if err != nil {
			logger.Error("Failed to create transaction", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
			return
		}

		// Update subscription
		subscription.Status = "Active"

		monthlyPrice := getMonthlyPrice(subscription)
		totalDue := subscription.DueAmount + monthlyPrice

		if transactionInfo.Amount >= totalDue {
			// Full payment or overpayment
			monthsPaid := int(transactionInfo.Amount / monthlyPrice)
			subscription.PaidUntil = addMonths(subscription.PaidUntil, monthsPaid)
			subscription.DueAmount = 0
		} else {
			// Partial payment
			subscription.DueAmount = totalDue - transactionInfo.Amount
		}

		// Set new RenewalDate
		subscription.RenewalDate = getFirstDayOfNextMonth(subscription.PaidUntil)

		err = h.subscriptionRepo.UpdateSubscription(c.Request.Context(), subscription)
		if err != nil {
			logger.Error("Failed to update subscription", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":      "Cash transaction processed successfully",
			"transaction":  transaction,
			"subscription": subscription,
		})
	}
}

func getMonthlyPrice(subscription *models.Subscription) float64 {
	price, _ := strconv.ParseFloat(subscription.PackagePrice, 64)
	return price
}

func addMonths(date time.Time, months int) time.Time {
	return date.AddDate(0, months, 0)
}
