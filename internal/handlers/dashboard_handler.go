package handler

import (
	service "gin-investment-tracker/internal/services"
	"gin-investment-tracker/internal/util"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	service service.DashboardServiceInterface
}

func NewDashboardHandler(svc service.DashboardServiceInterface) *DashboardHandler {
	return &DashboardHandler{service: svc}
}

func (h *DashboardHandler) SetRoutes(rg *gin.RouterGroup) {
	dashboard := rg.Group("/dashboard")
	dashboard.GET("", h.GetDashboardData)
}

// GetDashboardData godoc
// @Summary Get dashboard data
// @Description Returns quick insights, top holdings, and recent transactions for the authenticated user
// @Tags dashboard
// @Produce json
// @Success 200 {object} dto.DashboardDataDto
// @Failure 400 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/dashboard [get]
// @Security BearerAuth
func (h *DashboardHandler) GetDashboardData(c *gin.Context) {
	slog.Info("request started", "handler", "DashboardHandler.GetDashboardData", "method", c.Request.Method, "path", c.Request.URL.Path)

	userID, ok := util.GetUserIDFromContext(c)
	if !ok {
		slog.Warn("failed to parse user ID from context", "handler", "DashboardHandler.GetDashboardData")
		util.SendErrorResponse(c, http.StatusBadRequest, "error while parsing the userId")
		return
	}

	data, err := h.service.GetDashboardData(c.Request.Context(), userID)
	if err != nil {
		util.HandleError(c, err, "DashboardHandler.GetDashboardData")
		return
	}

	slog.Info("dashboard data retrieved", "handler", "DashboardHandler.GetDashboardData", "userID", userID)
	util.SendResponse(c, http.StatusOK, data)
}
