package handlers

import (
	"context"
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

type SubscriptionHandler struct {
	repo        repositories.SubscriptionRepository
	packageRepo repositories.PackageRepository
	deviceRepo  repositories.DeviceRepository
}

func NewSubscriptionHandler(
	repo repositories.SubscriptionRepository,
	packageRepo repositories.PackageRepository,
	deviceRepo repositories.DeviceRepository) *SubscriptionHandler {
	return &SubscriptionHandler{
		repo:        repo,
		packageRepo: packageRepo,
		deviceRepo:  deviceRepo,
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
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if subscription.CustomerID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "CustomerID is required"})
			return
		}

		// Validate package
		var pkg interface{}
		var err error
		if subscription.Type == models.Internet {
			pkg, err = h.packageRepo.GetInternetPackageByID(c.Request.Context(), subscription.PackageID)
		} else if subscription.Type == models.CableTV {
			pkg, err = h.packageRepo.GetCableTVPackageByID(c.Request.Context(), subscription.PackageID)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription type"})
			return
		}

		if err != nil {
			logger.Error("Failed to get package", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid package ID"})
			return
		}

		// Set package details
		if internetPkg, ok := pkg.(*models.InternetPackage); ok {
			subscription.PackageName = internetPkg.PackageName
			subscription.PackagePrice = internetPkg.Price
		} else if cableTVPkg, ok := pkg.(*models.CableTVPackage); ok {
			subscription.PackageName = cableTVPkg.PackageName
			subscription.PackagePrice = cableTVPkg.Price
		}

		// Generate a unique ID for the subscription
		subscription.ID = uuid.New().String()

		// Set StartDate to current time
		subscription.StartDate = time.Now()

		// Set RenewalDate to the first day of next month
		subscription.RenewalDate = getFirstDayOfNextMonth(subscription.StartDate)

		// Set PaidUntil to StartDate (assuming no payment has been made yet)
		subscription.PaidUntil = subscription.StartDate

		// Initialize DueAmount to 0
		subscription.DueAmount = 0

		err = h.repo.CreateSubscription(c.Request.Context(), &subscription)
		if err != nil {
			logger.Error("Failed to create subscription", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
			return
		}

		logger.Info("Subscription created successfully", zap.String("id", subscription.ID))
		c.JSON(http.StatusCreated, gin.H{"message": "Subscription created successfully", "subscription": subscription})
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
			logger.Error("Failed to bind JSON", zap.Error(err))
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

		// Update fields
		existingSubscription.Type = updateData.Type
		existingSubscription.PackageName = updateData.PackageName
		existingSubscription.PackagePrice = updateData.PackagePrice
		existingSubscription.MonthlyDiscount = updateData.MonthlyDiscount
		existingSubscription.Status = updateData.Status
		existingSubscription.DeviceID = updateData.DeviceID
		existingSubscription.DueAmount = updateData.DueAmount

		err = h.repo.UpdateSubscription(c.Request.Context(), existingSubscription)
		if err != nil {
			logger.Error("Failed to update subscription", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Subscription updated successfully", "subscription": existingSubscription})
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
func (h *SubscriptionHandler) ProcessExpiredSubscriptions() {
	ctx := context.Background()
	expiredSubscriptions, err := h.repo.GetExpiredSubscriptions(ctx)
	if err != nil {
		logger.Error("Failed to get expired subscriptions", zap.Error(err))
		return
	}

	for _, subscription := range expiredSubscriptions {
		if subscription.DeviceID != "" {
			err := h.deviceRepo.MarkDeviceForCollection(ctx, subscription.DeviceID)
			if err != nil {
				logger.Error("Failed to mark device for collection", zap.Error(err), zap.String("deviceID", subscription.DeviceID))
			}
		}
		subscription.Status = "Expired"
		err := h.repo.UpdateSubscription(ctx, &subscription)
		if err != nil {
			logger.Error("Failed to update subscription status", zap.Error(err), zap.String("subscriptionID", subscription.ID))
		}
	}
}
