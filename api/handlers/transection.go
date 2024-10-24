package handlers

import (
	"github.com/timam/uttarawave-backend/internals/models"
	"time"
)

func getMonthlyPrice(subscription *models.Subscription) float64 {
	return subscription.PackagePrice
}

func addMonths(date time.Time, months int) time.Time {
	return date.AddDate(0, months, 0)
}
