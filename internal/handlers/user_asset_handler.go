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

type UserAssetHandler struct {
	service service.UserAssetServiceInterface
}

func NewUserAssetHandler(svc service.UserAssetServiceInterface) *UserAssetHandler {
	return &UserAssetHandler{service: svc}
}

func (h *UserAssetHandler) SetRoutes(rg *gin.RouterGroup) {
	userAssets := rg.Group("/user-assets")
	userAssets.POST("", h.Create)
	userAssets.GET("", h.GetByUserID)
	userAssets.DELETE("/:userAssetId", h.Delete)
}

// Create godoc
// @Summary Assign asset to user
// @Description Link an asset to a user
// @Tags user-assets
// @Accept json
// @Produce json
// @Param payload body dto.CreateUserAssetRequest true "Asset assignment payload"
// @Success 201 {object} map[string]string
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/user-assets [post]
// @Security BearerAuth
func (h *UserAssetHandler) Create(c *gin.Context) {
	slog.Info("request started", "handler", "UserAssetHandler.Create", "method", c.Request.Method, "path", c.Request.URL.Path)

	userID, ok := util.GetUserIDFromContext(c)
	if !ok {
		slog.Warn("failed to parse user ID from context", "handler", "UserAssetHandler.Create")
		util.SendErrorResponse(c, http.StatusBadRequest, "error while parsing the userId")
		return
	}

	var req dto.CreateUserAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			slog.Warn("validation failed", "handler", "UserAssetHandler.Create", "userID", userID, "error", ve)
			util.SendErrorResponse(c, http.StatusBadRequest, util.FormatValidationErrors(err))
			return
		}
		slog.Warn("failed to bind request body", "handler", "UserAssetHandler.Create", "userID", userID, "error", err)
		util.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userAsset, err := h.service.Create(c.Request.Context(), userID, &req)
	if err != nil {
		util.HandleError(c, err, "UserAssetHandler.Create")
		return
	}

	slog.Info("user asset created", "handler", "UserAssetHandler.Create", "userID", userAsset.UserID, "assetID", userAsset.AssetID)
	util.SendResponse(c, http.StatusCreated, map[string]any{
		"message":    fmt.Sprintf("asset %d assigned to user %d", userAsset.AssetID, userAsset.UserID),
		"user_asset": userAsset,
	})
}

// GetByUserID godoc
// @Summary List assets for a user
// @Description Get all asset assignments for a user with pagination
// @Tags user-assets
// @Produce json
// @Param limit query int false "Number of records to return (default: 50, max: 200)"
// @Param offset query int false "Number of records to skip (default: 0)"
// @Success 200 {array} model.UserAsset
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/user-assets [get]
// @Security BearerAuth
func (h *UserAssetHandler) GetByUserID(c *gin.Context) {
	slog.Info("request started", "handler", "UserAssetHandler.GetByUserID", "method", c.Request.Method, "path", c.Request.URL.Path)

	userID, ok := util.GetUserIDFromContext(c)
	if !ok {
		slog.Warn("failed to parse user ID from context", "handler", "UserAssetHandler.GetByUserID")
		util.SendErrorResponse(c, http.StatusBadRequest, "error while parsing the userId")
		return
	}

	limit, offset, err := parsePaginationParams(c)
	if err != nil {
		slog.Warn("invalid pagination params", "handler", "UserAssetHandler.GetByUserID", "userID", userID, "error", err)
		util.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userAssets, err := h.service.GetByUserID(c.Request.Context(), userID, limit, offset)
	if err != nil {
		util.HandleError(c, err, "UserAssetHandler.GetByUserID")
		return
	}

	slog.Info("user assets retrieved", "handler", "UserAssetHandler.GetByUserID", "userID", userID, "count", len(userAssets), "limit", limit, "offset", offset)
	util.SendResponse(c, http.StatusOK, userAssets)
}

// Delete godoc
// @Summary Remove asset assignment
// @Description Delete a user-asset link by its ID
// @Tags user-assets
// @Produce json
// @Param userAssetId path int true "User Asset ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/user-assets/{userAssetId} [delete]
// @Security BearerAuth
func (h *UserAssetHandler) Delete(c *gin.Context) {
	slog.Info("request started", "handler", "UserAssetHandler.Delete", "method", c.Request.Method, "path", c.Request.URL.Path)

	userID, ok := util.GetUserIDFromContext(c)
	if !ok {
		slog.Warn("failed to parse user ID from context", "handler", "UserAssetHandler.Delete")
		util.SendErrorResponse(c, http.StatusBadRequest, "error while parsing the userId")
		return
	}

	userAssetID, err := parseIntegerID(c, "userAssetId")
	if err != nil {
		slog.Warn("invalid user asset ID", "handler", "UserAssetHandler.Delete", "userID", userID, "error", err)
		util.SendErrorResponse(c, http.StatusBadRequest, "invalid user asset id")
		return
	}

	if err := h.service.Delete(c.Request.Context(), userID, userAssetID); err != nil {
		util.HandleError(c, err, "UserAssetHandler.Delete")
		return
	}

	slog.Info("user asset deleted", "handler", "UserAssetHandler.Delete", "userID", userID, "userAssetID", userAssetID)
	util.SendResponse(c, http.StatusOK, map[string]string{"message": "user asset deleted successfully"})
}
