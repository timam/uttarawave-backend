package handlers

import (
	"github.com/timam/uttarawave-backend/models"
	"strconv"
	"time"
)

func getMonthlyPrice(subscription *models.Subscription) float64 {
	price, _ := strconv.ParseFloat(subscription.PackagePrice, 64)
	return price
}

func addMonths(date time.Time, months int) time.Time {
	return date.AddDate(0, months, 0)
}
