package handler

import (
	"fmt"
	service "gin-investment-tracker/internal/services"
	"gin-investment-tracker/internal/util"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatementHandler struct {
	service service.StatementServiceInterface
}

func NewCasStatementHandler(svc service.StatementServiceInterface) *StatementHandler {
	return &StatementHandler{service: svc}
}

func (h *StatementHandler) SetRoutes(rg *gin.RouterGroup) {
	casStatement := rg.Group("cas-statement")
	casStatement.POST("", h.ProcessCasStatement)
}

// ProcessCasStatment godoc
// @Summary Process CAS statement
// @Description Upload a CAS (Consolidated Account Statement) PDF file to import mutual fund transactions. Processing happens asynchronously in the background.
// @Tags cas-statement
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CAS statement PDF file"
// @Param password formData string false "PDF password (if password protected)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorBody
// @Failure 401 {object} util.ErrorBody
// @Router /api/cas-statement [post]
// @Security BearerAuth
func (h *StatementHandler) ProcessCasStatement(c *gin.Context) {
	filePassword := c.PostForm("password")

	file, err := c.FormFile("file")
	if err != nil {
		slog.Warn("No cas statement provided", "handler", "StatementHandler.ProcessCasStatement")
		util.SendErrorResponse(c, http.StatusBadRequest, "No cas statement provided")
		return
	}

	userID, ok := util.GetUserIDFromContext(c)
	if !ok {
		slog.Warn("failed to parse user ID from context", "handler", "StatementHandler.ProcessCasStatement")
		util.SendErrorResponse(c, http.StatusBadRequest, "error while parsing the userId")
		return
	}

	file.Filename = fmt.Sprintf("cas_%d", userID)
	h.service.ProcessCasFile(c.Request.Context(), file, filePassword, userID)

	util.SendResponse(c, http.StatusAccepted, map[string]string{
		"message": "Your file is being processed. Please check back in a few minutes.",
	})
}
