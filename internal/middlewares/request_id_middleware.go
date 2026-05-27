package middleware

import (
	"gin-investment-tracker/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Set(util.RequestIDKey, requestID)
		c.Header("X-Request-ID", requestID)

		// Attaching request id to logger
		log := util.Logger.With("request_id", requestID)
		c.Request = c.Request.WithContext(util.WithLogger(c.Request.Context(), log))
		c.Next()
	}
}
