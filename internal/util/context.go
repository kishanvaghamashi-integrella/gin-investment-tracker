package util

import "github.com/gin-gonic/gin"

const (
	UserIDKey    = "user_id"
	RequestIDKey = "request_id"
)

// GetUserIDFromContext retrieves the user ID from the Gin context
func GetUserIDFromContext(c *gin.Context) (int64, bool) {
	val, exists := c.Get(UserIDKey)
	if !exists {
		return 0, false
	}
	userID, ok := val.(int64)
	return userID, ok
}

// GetRequestIDFromContext retrieves the request ID from the Gin context
func GetRequestIDFromContext(c *gin.Context) string {
	val, exists := c.Get(RequestIDKey)
	if !exists {
		return ""
	}
	id, ok := val.(string)
	if !ok {
		return ""
	}
	return id
}
