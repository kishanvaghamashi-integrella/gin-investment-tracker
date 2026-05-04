package middleware

import (
	"gin-investment-tracker/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwt_token")
		if err != nil {
			util.SendErrorResponse(c, http.StatusUnauthorized, "missing auth cookie")
			c.Abort()
			return
		}

		userID, err := util.ValidateToken(tokenString)
		if err != nil {
			util.SendErrorResponse(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		c.Set(util.UserIDKey, userID)
		c.Next()
	}
}
