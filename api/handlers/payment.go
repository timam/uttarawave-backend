package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	models2 "github.com/timam/uttarawave-backend/internals/models"
	repositories2 "github.com/timam/uttarawave-backend/internals/repositories"
	"github.com/timam/uttarawave-backend/pkg/logger"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type PaymentHandler struct {
	paymentRepo      repositories2.PaymentRepository
	subscriptionRepo repositories2.SubscriptionRepository
	invoiceRepo      repositories2.InvoiceRepository
}

func NewPaymentHandler(pr repositories2.PaymentRepository, sr repositories2.SubscriptionRepository, ir repositories2.InvoiceRepository) *PaymentHandler {
	return &PaymentHandler{
		paymentRepo:      pr,
		subscriptionRepo: sr,
		invoiceRepo:      ir,
	}
}

func (h *PaymentHandler) CreatePayment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payment models2.Payment
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

			if invoice.Status == models2.InvoicePaid {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invoice is already paid"})
				return
			}

			if payment.Amount >= invoice.Amount {
				invoice.Status = models2.InvoicePaid
				invoice.PaidDate = &payment.PaidAt
			} else {
				invoice.Status = models2.InvoicePending
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
