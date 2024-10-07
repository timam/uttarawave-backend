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

type subscriptionHandler struct {
	repo        repositories.SubscriptionRepository
	packageRepo repositories.PackageRepository
}

func NewSubscriptionHandler(repo repositories.SubscriptionRepository, packageRepo repositories.PackageRepository) *subscriptionHandler {
	return &subscriptionHandler{
		repo:        repo,
		packageRepo: packageRepo,
	}
}

func (h *subscriptionHandler) CreateSubscription() gin.HandlerFunc {
	return func(c *gin.Context) {
		var subscription models.Subscription

		if err := c.ShouldBindJSON(&subscription); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func getFirstDayOfNextMonth(date time.Time) time.Time {
	firstOfMonth := date.AddDate(0, 1, 0)
	return time.Date(firstOfMonth.Year(), firstOfMonth.Month(), 1, 0, 0, 0, 0, firstOfMonth.Location())
}
func (h *subscriptionHandler) GetSubscription() gin.HandlerFunc {
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

func (h *subscriptionHandler) UpdateSubscription() gin.HandlerFunc {
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
		existingSubscription.Discount = updateData.Discount
		existingSubscription.Status = updateData.Status
		existingSubscription.DeviceID = updateData.DeviceID

		err = h.repo.UpdateSubscription(c.Request.Context(), existingSubscription)
		if err != nil {
			logger.Error("Failed to update subscription", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Subscription updated successfully", "subscription": existingSubscription})
	}
}

func (h *subscriptionHandler) DeleteSubscription() gin.HandlerFunc {
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

func (h *subscriptionHandler) GetAllSubscriptions() gin.HandlerFunc {
	return func(c *gin.Context) {
		subscriptions, err := h.repo.GetAllSubscriptions(c.Request.Context())
		if err != nil {
			logger.Error("Failed to get all subscriptions", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all subscriptions"})
			return
		}

		c.JSON(http.StatusOK, subscriptions)
	}
}
