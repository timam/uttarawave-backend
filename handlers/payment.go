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

type PaymentHandler struct {
	paymentRepo      repositories.PaymentRepository
	subscriptionRepo repositories.SubscriptionRepository
	invoiceRepo      repositories.InvoiceRepository
}

func NewPaymentHandler(pr repositories.PaymentRepository, sr repositories.SubscriptionRepository, ir repositories.InvoiceRepository) *PaymentHandler {
	return &PaymentHandler{
		paymentRepo:      pr,
		subscriptionRepo: sr,
		invoiceRepo:      ir,
	}
}

func (h *PaymentHandler) CreatePayment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payment models.Payment
		if err := c.ShouldBindJSON(&payment); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		payment.ID = uuid.New().String()
		payment.PaidAt = time.Now()

		if payment.InvoiceID != nil {
			invoice, err := h.invoiceRepo.GetInvoiceByID(c.Request.Context(), *payment.InvoiceID)
			if err != nil {
				logger.Error("Failed to get invoice", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get invoice"})
				return
			}

			if invoice.Status == models.InvoicePaid {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invoice is already paid"})
				return
			}

			if payment.Amount >= invoice.Amount {
				invoice.Status = models.InvoicePaid
				invoice.PaidDate = &payment.PaidAt
			} else {
				invoice.Status = models.InvoicePending
			}

			err = h.invoiceRepo.UpdateInvoice(c.Request.Context(), invoice)
			if err != nil {
				logger.Error("Failed to update invoice", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update invoice"})
				return
			}

			if invoice.SubscriptionID != nil {
				subscription, err := h.subscriptionRepo.GetSubscription(c.Request.Context(), *invoice.SubscriptionID)
				if err != nil {
					logger.Error("Failed to get subscription", zap.Error(err))
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subscription"})
					return
				}

				dueAmount, _ := strconv.ParseFloat(subscription.DueAmount, 64)
				dueAmount -= payment.Amount
				subscription.DueAmount = strconv.FormatFloat(dueAmount, 'f', 2, 64)

				if dueAmount <= 0 {
					subscription.Status = "Active"
					subscription.PaidUntil = addMonths(subscription.PaidUntil, 1)
					subscription.RenewalDate = getFirstDayOfNextMonth(subscription.PaidUntil)
				}

				err = h.subscriptionRepo.UpdateSubscription(c.Request.Context(), subscription)
				if err != nil {
					logger.Error("Failed to update subscription", zap.Error(err))
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription"})
					return
				}
			}
		}

		err := h.paymentRepo.CreatePayment(c.Request.Context(), &payment)
		if err != nil {
			logger.Error("Failed to create payment", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment"})
			return
		}

		c.JSON(http.StatusCreated, payment)
	}
}
