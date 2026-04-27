package handler

import (
	"fmt"
	service "gin-investment-tracker/internal/services"
	"gin-investment-tracker/internal/util"
	"log"
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
	casStatement.POST("", h.ProcessCasStatment)
}

func (h *StatementHandler) ProcessCasStatment(c *gin.Context) {
	filePassword := c.PostForm("password")

	file, err := c.FormFile("file")
	if err != nil {
		slog.Warn("No cas statement provided", "handler", "StatementHandler.ProcessCasStatment")
		util.SendErrorResponse(c, http.StatusBadRequest, "No cas statement provided")
		return
	}

	userID, ok := util.GetUserIDFromContext(c)
	if !ok {
		slog.Warn("failed to parse user ID from context", "handler", "StatementHandler.ProcessCasStatment")
		util.SendErrorResponse(c, http.StatusBadRequest, "error while parsing the userId")
		return
	}

	file.Filename = fmt.Sprintf("cas_%d", userID)

	h.service.ProcessCasFile(c.Request.Context(), file, filePassword, userID)

	log.Printf("file name is %s and password is %s - %d", file.Filename, filePassword, userID)
	util.SendResponse(c, 200, map[string]string{
		"messege": "Your file has been processing, please check after few minutes",
	})
}
