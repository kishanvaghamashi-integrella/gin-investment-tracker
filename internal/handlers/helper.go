package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func parseIntegerID(c *gin.Context, param string) (int64, error) {
	return strconv.ParseInt(c.Param(param), 10, 64)
}

func parsePaginationParams(c *gin.Context) (int, int, error) {
	const (
		defaultLimit = 50
		maxLimit     = 200
	)

	limit := defaultLimit
	offset := 0

	if limitValue := c.Query("limit"); limitValue != "" {
		parsedLimit, err := strconv.Atoi(limitValue)
		if err != nil || parsedLimit <= 0 {
			return 0, 0, fmt.Errorf("limit must be a positive integer")
		}
		if parsedLimit > maxLimit {
			parsedLimit = maxLimit
		}
		limit = parsedLimit
	}

	if offsetValue := c.Query("offset"); offsetValue != "" {
		parsedOffset, err := strconv.Atoi(offsetValue)
		if err != nil || parsedOffset < 0 {
			return 0, 0, fmt.Errorf("offset must be a non-negative integer")
		}
		offset = parsedOffset
	}

	return limit, offset, nil
}
