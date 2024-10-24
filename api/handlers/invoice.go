package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/timam/uttarawave-backend/internals/models"
	repositories2 "github.com/timam/uttarawave-backend/internals/repositories"
	"github.com/timam/uttarawave-backend/pkg/logger"
	"go.uber.org/zap"
	"net/http"
)

type InvoiceHandler struct {
	invoiceRepo      repositories2.InvoiceRepository
	subscriptionRepo repositories2.SubscriptionRepository
}

func NewInvoiceHandler(ir repositories2.InvoiceRepository, sr repositories2.SubscriptionRepository) *InvoiceHandler {
	return &InvoiceHandler{
		invoiceRepo:      ir,
		subscriptionRepo: sr,
	}
}

func (h *InvoiceHandler) CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var invoice models.Invoice
		if err := c.ShouldBindJSON(&invoice); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		invoice.ID = uuid.New().String()
		invoice.Status = models.InvoicePending

		if invoice.SubscriptionID != nil {
			subscription, err := h.subscriptionRepo.GetSubscription(c.Request.Context(), *invoice.SubscriptionID)
			if err != nil {
				logger.Error("Failed to get subscription", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subscription"})
				return
			}

			invoice.CustomerID = subscription.CustomerID
			invoice.Amount = getMonthlyPrice(subscription)
			invoice.DueDate = subscription.RenewalDate
		}

		err := h.invoiceRepo.CreateInvoice(c.Request.Context(), &invoice)
		if err != nil {
			logger.Error("Failed to create invoice", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invoice"})
			return
		}

		c.JSON(http.StatusCreated, invoice)
	}
}

// Add other invoice handler methods (GetInvoice, UpdateInvoice, GetAllInvoices) here
