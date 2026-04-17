package handler

import (
	"errors"
	"fmt"
	dto "gin-investment-tracker/internal/dtos"
	service "gin-investment-tracker/internal/services"
	"gin-investment-tracker/internal/util"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TransactionHandler struct {
	service service.TransactionServiceInterface
}

func NewTransactionHandler(svc service.TransactionServiceInterface) *TransactionHandler {
	return &TransactionHandler{service: svc}
}

func (h *TransactionHandler) SetRoutes(rg *gin.RouterGroup) {
	transactions := rg.Group("/transactions")
	transactions.POST("", h.Create)
	transactions.GET("", h.GetAllByUserID)
	transactions.PUT("/:txnId", h.Update)
	transactions.DELETE("/:txnId", h.Delete)
}

// Create godoc
// @Summary Create transaction
// @Description Create a new transaction for a user asset
// @Tags transactions
// @Accept json
// @Produce json
// @Param payload body dto.CreateTransactionRequest true "Create transaction payload"
// @Success 201 {object} map[string]any
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/transactions [post]
// @Security BearerAuth
func (h *TransactionHandler) Create(c *gin.Context) {
	slog.Info("request started", "handler", "TransactionHandler.Create", "method", c.Request.Method, "path", c.Request.URL.Path)

	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			slog.Warn("validation failed", "handler", "TransactionHandler.Create", "error", ve)
			util.SendErrorResponse(c, http.StatusBadRequest, util.FormatValidationErrors(err))
			return
		}
		slog.Warn("failed to bind request body", "handler", "TransactionHandler.Create", "error", err)
		util.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, ok := util.GetUserIDFromContext(c)
	if !ok {
		slog.Warn("failed to parse user ID from context", "handler", "TransactionHandler.Create")
		util.SendErrorResponse(c, http.StatusBadRequest, "error while parsing the userId")
		return
	}

	txn, err := h.service.Create(c.Request.Context(), &req, userID)
	if err != nil {
		util.HandleError(c, err, "TransactionHandler.Create")
		return
	}

	slog.Info("transaction created", "handler", "TransactionHandler.Create", "transactionID", txn.ID)
	util.SendResponse(c, http.StatusCreated, map[string]any{
		"message":     fmt.Sprintf("transaction created with id %d", txn.ID),
		"transaction": txn,
	})
}

// GetAllByUserID godoc
// @Summary List transactions for the authenticated user
// @Description Get all transactions for the current user with pagination
// @Tags transactions
// @Produce json
// @Param limit query int false "Number of records to return (default: 50, max: 200)"
// @Param offset query int false "Number of records to skip (default: 0)"
// @Success 200 {array} dto.ResponseTransactionDto
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/transactions [get]
// @Security BearerAuth
func (h *TransactionHandler) GetAllByUserID(c *gin.Context) {
	slog.Info("request started", "handler", "TransactionHandler.GetAllByUserID", "method", c.Request.Method, "path", c.Request.URL.Path)

	userID, ok := util.GetUserIDFromContext(c)
	if !ok {
		slog.Warn("failed to parse user ID from context", "handler", "TransactionHandler.GetAllByUserID")
		util.SendErrorResponse(c, http.StatusBadRequest, "error while parsing the userId")
		return
	}

	limit, offset, err := parsePaginationParams(c)
	if err != nil {
		slog.Warn("invalid pagination params", "handler", "TransactionHandler.GetAllByUserID", "error", err)
		util.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	transactions, err := h.service.GetAllByUserID(c.Request.Context(), userID, limit, offset)
	if err != nil {
		util.HandleError(c, err, "TransactionHandler.GetAllByUserID")
		return
	}

	slog.Info("transactions retrieved", "handler", "TransactionHandler.GetAllByUserID", "userID", userID, "count", len(transactions))
	util.SendResponse(c, http.StatusOK, transactions)
}

// Update godoc
// @Summary Update a transaction
// @Description Update an existing transaction by ID
// @Tags transactions
// @Accept json
// @Produce json
// @Param txnId path int true "Transaction ID"
// @Param payload body dto.UpdateTransactionRequest true "Update transaction payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/transactions/{txnId} [put]
// @Security BearerAuth
func (h *TransactionHandler) Update(c *gin.Context) {
	slog.Info("request started", "handler", "TransactionHandler.Update", "method", c.Request.Method, "path", c.Request.URL.Path)

	id, err := parseIntegerID(c, "txnId")
	if err != nil {
		slog.Warn("invalid transaction ID", "handler", "TransactionHandler.Update", "txnId", c.Param("txnId"))
		util.SendErrorResponse(c, http.StatusBadRequest, "invalid transaction id")
		return
	}

	var req dto.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			slog.Warn("validation failed", "handler", "TransactionHandler.Update", "error", ve)
			util.SendErrorResponse(c, http.StatusBadRequest, util.FormatValidationErrors(err))
			return
		}
		slog.Warn("failed to bind request body", "handler", "TransactionHandler.Update", "error", err)
		util.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Update(c.Request.Context(), id, &req); err != nil {
		util.HandleError(c, err, "TransactionHandler.Update")
		return
	}

	slog.Info("transaction updated", "handler", "TransactionHandler.Update", "transactionID", id)
	util.SendResponse(c, http.StatusOK, map[string]string{"message": "transaction updated successfully"})
}

// Delete godoc
// @Summary Delete a transaction
// @Description Delete a transaction by ID
// @Tags transactions
// @Produce json
// @Param txnId path int true "Transaction ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/transactions/{txnId} [delete]
// @Security BearerAuth
func (h *TransactionHandler) Delete(c *gin.Context) {
	slog.Info("request started", "handler", "TransactionHandler.Delete", "method", c.Request.Method, "path", c.Request.URL.Path)

	id, err := parseIntegerID(c, "txnId")
	if err != nil {
		slog.Warn("invalid transaction ID", "handler", "TransactionHandler.Delete", "txnId", c.Param("txnId"))
		util.SendErrorResponse(c, http.StatusBadRequest, "invalid transaction id")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		util.HandleError(c, err, "TransactionHandler.Delete")
		return
	}

	slog.Info("transaction deleted", "handler", "TransactionHandler.Delete", "transactionID", id)
	util.SendResponse(c, http.StatusOK, map[string]string{"message": "transaction deleted successfully"})
}

