package middlewares

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/timam/uttarawave-backend/pkg/metrics"
	"strconv"
	"time"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Get the request size
		requestSize := float64(c.Request.ContentLength)

		// Create a custom ResponseWriter to capture the response size
		w := &responseWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		status := c.Writer.Status()
		duration := time.Since(start)
		responseSize := float64(w.body.Len())

		metrics.RecordHTTPMetrics(c.Request.Method, c.FullPath(), strconv.Itoa(status), duration.Seconds(), requestSize, responseSize)
	}
}
