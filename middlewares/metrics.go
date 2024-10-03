package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/timam/uttarawave-finance-backend/pkg/metrics"
	"strconv"
	"time"
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := c.Writer.Status()
		duration := time.Since(start)

		metrics.RecordHTTPMetrics(c.Request.Method, c.FullPath(), strconv.Itoa(status), duration.Seconds())
	}
}
