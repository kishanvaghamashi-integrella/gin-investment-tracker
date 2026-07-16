package metrics

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Increment active requests
		httpRequestsActive.Inc()
		defer httpRequestsActive.Dec()

		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()

		// Record metrics
		endpoint := c.FullPath()
		method := c.Request.Method
		status := fmt.Sprintf("%d", c.Writer.Status())

		// Record counter
		httpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()

		// Record histogram (duration)
		httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
	}
}
