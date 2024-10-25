package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/timam/uttarawave-backend/internals/models"
	repositories "github.com/timam/uttarawave-backend/internals/repositories"
	"github.com/timam/uttarawave-backend/pkg/logger"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type SubscriptionHandler struct {
	repo        repositories.SubscriptionRepository
	packageRepo repositories.PackageRepository
	deviceRepo  repositories.DeviceRepository
	invoiceRepo repositories.InvoiceRepository
}

func NewSubscriptionHandler(
	repo repositories.SubscriptionRepository,
	packageRepo repositories.PackageRepository,
	deviceRepo repositories.DeviceRepository,
	invoiceRepo repositories.InvoiceRepository) *SubscriptionHandler {
	return &SubscriptionHandler{
		repo:        repo,
		packageRepo: packageRepo,
		deviceRepo:  deviceRepo,
		invoiceRepo: invoiceRepo,
	}
}

func getFirstDayOfNextMonth(date time.Time) time.Time {
	currentYear, currentMonth, _ := date.Date()
	firstOfNextMonth := time.Date(currentYear, currentMonth+1, 1, 0, 0, 0, 0, date.Location())
	return firstOfNextMonth
}

func (h *SubscriptionHandler) CreateSubscription() gin.HandlerFunc {
	return func(c *gin.Context) {
		var subscription models.Subscription
		if err := c.ShouldBindJSON(&subscription); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate package
		pkg, err := h.packageRepo.GetPackageByID(c.Request.Context(), subscription.PackageID)
		if err != nil {
			logger.Error("Failed to get package", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid package ID"})
			return
		}

		// Set subscription details
		subscription.ID = uuid.New().String()
		subscription.PackagePrice = pkg.Price
		subscription.Status = "Active"
		subscription.StartDate = time.Now()
		subscription.RenewalDate = getFirstDayOfNextMonth(subscription.StartDate)
		subscription.PaidUntil = subscription.StartDate
		subscription.DueAmount = strconv.FormatFloat(pkg.Price, 'f', 2, 64)

		err = h.repo.CreateSubscription(c.Request.Context(), &subscription)
		if err != nil {
			logger.Error("Failed to create subscription", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
			return
		}

		c.JSON(http.StatusCreated, subscription)
	}
}
func (h *SubscriptionHandler) GetSubscription() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Subscription ID is required"})
			return
		}

		subscription, err := h.repo.GetSubscription(c.Request.Context(), id)
		if err != nil {
			logger.Error("Failed to get subscription", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subscription"})
			return
		}

		if subscription == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}

		c.JSON(http.StatusOK, subscription)
	}
}

func (h *SubscriptionHandler) UpdateSubscription() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Subscription ID is required"})
			return
		}

		var updateData models.Subscription
		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		existingSubscription, err := h.repo.GetSubscription(c.Request.Context(), id)
		if err != nil {
			logger.Error("Failed to get subscription", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subscription"})
			return
		}

		if existingSubscription == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}

		// Update fields based on provided data
		if updateData.Status != "" {
			existingSubscription.Status = updateData.Status
		}
		if updateData.PackageID != "" {
			existingSubscription.PackageID = updateData.PackageID
			// Fetch the new package and update related fields
			newPackage, err := h.packageRepo.GetPackageByID(c.Request.Context(), updateData.PackageID)
			if err == nil {
				existingSubscription.PackagePrice = newPackage.Price
			}
		}
		if updateData.DueAmount != "" {
			existingSubscription.DueAmount = updateData.DueAmount
		}
		// Update other fields as needed

		err = h.repo.UpdateSubscription(c.Request.Context(), existingSubscription)
		if err != nil {
			logger.Error("Failed to update subscription data", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription data"})
			return
		}

		c.JSON(http.StatusOK, existingSubscription)
	}
}

func (h *SubscriptionHandler) DeleteSubscription() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Subscription ID is required"})
			return
		}

		err := h.repo.DeleteSubscription(c.Request.Context(), id)
		if err != nil {
			logger.Error("Failed to delete subscription", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subscription"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Subscription deleted successfully"})
	}
}

func (h *SubscriptionHandler) GetAllSubscriptions() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))

		subscriptions, totalCount, err := h.repo.GetSubscriptionsPaginated(c.Request.Context(), page, pageSize)
		if err != nil {
			logger.Error("Failed to get subscriptions", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subscriptions"})
			return
		}

		response := gin.H{
			"subscriptions": subscriptions,
			"totalCount":    totalCount,
			"page":          page,
			"pageSize":      pageSize,
		}

		logger.Info("Retrieved subscriptions", zap.Int("count", len(subscriptions)), zap.Int("page", page), zap.Int("pageSize", pageSize))
		c.JSON(http.StatusOK, response)
	}
}

//func (h *SubscriptionHandler) ProcessExpiredSubscriptions() {
//	ctx := context.Background()
//	expiredSubscriptions, err := h.repo.GetExpiredSubscriptions(ctx)
//	if err != nil {
//		logger.Error("Failed to get expired subscriptions", zap.Error(err))
//		return
//	}
//
//	for _, subscription := range expiredSubscriptions {
//		if subscription.DeviceID != "" {
//			err := h.deviceRepo.MarkDeviceForCollection(ctx, subscription.DeviceID)
//			if err != nil {
//				logger.Error("Failed to mark device for collection", zap.Error(err), zap.String("deviceID", subscription.DeviceID))
//			}
//		}
//		subscription.Status = "Expired"
//		err := h.repo.UpdateSubscription(ctx, &subscription)
//		if err != nil {
//			logger.Error("Failed to update subscription status", zap.Error(err), zap.String("subscriptionID", subscription.ID))
//		}
//	}
//}
