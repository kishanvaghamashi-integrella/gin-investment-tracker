package handler

import (
	service "gin-investment-tracker/internal/services"
	"gin-investment-tracker/internal/util"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type HoldingHandler struct {
	service service.HoldingServiceInterface
}

func NewHoldingHandler(svc service.HoldingServiceInterface) *HoldingHandler {
	return &HoldingHandler{service: svc}
}

func (h *HoldingHandler) SetRoutes(rg *gin.RouterGroup) {
	holdings := rg.Group("/holdings")
	holdings.GET("", h.Get)
}

// Get godoc
// @Summary List holdings for a user
// @Description Get all holdings for the authenticated user with pagination
// @Tags holdings
// @Produce json
// @Param limit query int false "Number of records to return (default: 50, max: 200)"
// @Param offset query int false "Number of records to skip (default: 0)"
// @Param sortby query string false "Sorting parameter (default: id)"
// @Param asset-name query string false "Name of asset (default: "")"
// @Success 200 {array} dto.HoldingResponseDto
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/holdings [get]
// @Security CookieAuth
func (h *HoldingHandler) Get(c *gin.Context) {
	slog.Info("request started", "handler", "HoldingHandler.Get", "method", c.Request.Method, "path", c.Request.URL.Path)

	userID, ok := util.GetUserIDFromContext(c)
	if !ok {
		slog.Warn("failed to parse user ID from context", "handler", "HoldingHandler.Get")
		util.SendErrorResponse(c, http.StatusBadRequest, "error while parsing the userId")
		return
	}

	limit, offset, err := parsePaginationParams(c)
	if err != nil {
		slog.Warn("invalid pagination params", "handler", "HoldingHandler.Get", "userID", userID, "error", err)
		util.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	sortByQuery := c.DefaultQuery("sortby", "id")
	assetNameQuery := strings.ToLower(c.DefaultQuery("asset-name", ""))

	holdings, err := h.service.GetAllByUserID(c.Request.Context(), userID, limit, offset, sortByQuery, assetNameQuery)
	if err != nil {
		util.HandleError(c, err, "HoldingHandler.Get")
		return
	}

	slog.Info("holdings retrieved", "handler", "HoldingHandler.Get", "userID", userID, "count", len(holdings), "limit", limit, "offset", offset)
	util.SendResponse(c, http.StatusOK, holdings)
}
