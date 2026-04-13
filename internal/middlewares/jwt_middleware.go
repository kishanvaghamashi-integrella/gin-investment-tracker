package middleware

import (
	"gin-investment-tracker/internal/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			util.SendErrorResponse(c, http.StatusUnauthorized, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			util.SendErrorResponse(c, http.StatusUnauthorized, "invalid authorization header format")
			c.Abort()
			return
		}

		userID, err := util.ValidateToken(parts[1])
		if err != nil {
			util.SendErrorResponse(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		c.Set(util.UserIDKey, userID)
		c.Next()
	}
}
