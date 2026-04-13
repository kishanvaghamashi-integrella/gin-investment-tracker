package util

import "github.com/gin-gonic/gin"

type ErrorBody struct {
	Error any `json:"error"`
}

func SendResponse(c *gin.Context, statusCode int, responseData any) {
	c.JSON(statusCode, responseData)
}

func SendErrorResponse(c *gin.Context, statusCode int, errData any) {
	c.JSON(statusCode, ErrorBody{Error: errData})
}
