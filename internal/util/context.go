package util

import "github.com/gin-gonic/gin"

const UserIDKey = "user_id"

// GetUserIDFromContext retrieves the user ID from the Gin context
func GetUserIDFromContext(c *gin.Context) (int64, bool) {
	val, exists := c.Get(UserIDKey)
	if !exists {
		return 0, false
	}
	userID, ok := val.(int64)
	return userID, ok
}
