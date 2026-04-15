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

type AssetHandler struct {
	service service.AssetServiceInterface
}

func NewAssetHandler(svc service.AssetServiceInterface) *AssetHandler {
	return &AssetHandler{service: svc}
}

// Create godoc
// @Summary Create asset
// @Description Create a new asset
// @Tags assets
// @Accept json
// @Produce json
// @Param payload body dto.CreateAssetRequest true "Create asset payload"
// @Success 201 {object} map[string]string
// @Failure 400 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/assets/ [post]
// @Security BearerAuth
func (h *AssetHandler) Create(c *gin.Context) {
	slog.Info("request started", "handler", "AssetHandler.Create", "method", c.Request.Method, "path", c.Request.URL.Path)

	var req dto.CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			slog.Warn("validation failed", "handler", "AssetHandler.Create", "error", ve)
			util.SendErrorResponse(c, http.StatusBadRequest, util.FormatValidationErrors(err))
			return
		}
		slog.Warn("failed to bind request body", "handler", "AssetHandler.Create", "error", err)
		util.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	asset, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		util.HandleError(c, err, "AssetHandler.Create")
		return
	}

	slog.Info("asset created", "handler", "AssetHandler.Create", "assetID", asset.ID)
	util.SendResponse(c, http.StatusCreated, map[string]any{
		"message": fmt.Sprintf("asset created with id %d", asset.ID),
		"asset":   asset,
	})
}

// GetByID godoc
// @Summary Get asset by ID
// @Description Retrieve a single asset by its ID
// @Tags assets
// @Produce json
// @Param assetId path int64 true "Asset ID"
// @Success 200 {object} model.Asset
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/assets/{assetId} [get]
// @Security BearerAuth
func (h *AssetHandler) GetByID(c *gin.Context) {
	slog.Info("request started", "handler", "AssetHandler.GetByID", "method", c.Request.Method, "path", c.Request.URL.Path)

	id, err := parseIntegerID(c, "assetId")
	if err != nil {
		slog.Warn("invalid asset ID", "handler", "AssetHandler.GetByID", "assetId", c.Param("assetId"))
		util.SendErrorResponse(c, http.StatusBadRequest, "invalid asset id")
		return
	}

	asset, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		util.HandleError(c, err, "AssetHandler.GetByID")
		return
	}

	slog.Info("asset retrieved", "handler", "AssetHandler.GetByID", "assetID", id)
	util.SendResponse(c, http.StatusOK, asset)
}

// GetAll godoc
// @Summary List assets
// @Description Retrieve assets with pagination
// @Tags assets
// @Produce json
// @Param limit query int false "Number of records to return (default: 50, max: 200)"
// @Param offset query int false "Number of records to skip (default: 0)"
// @Success 200 {array} model.Asset
// @Failure 400 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/assets/ [get]
// @Security BearerAuth
func (h *AssetHandler) GetAll(c *gin.Context) {
	slog.Info("request started", "handler", "AssetHandler.GetAll", "method", c.Request.Method, "path", c.Request.URL.Path)

	limit, offset, err := parsePaginationParams(c)
	if err != nil {
		slog.Warn("invalid pagination params", "handler", "AssetHandler.GetAll", "error", err)
		util.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	assets, err := h.service.GetAll(c.Request.Context(), limit, offset)
	if err != nil {
		util.HandleError(c, err, "AssetHandler.GetAll")
		return
	}

	slog.Info("assets retrieved", "handler", "AssetHandler.GetAll", "count", len(assets))
	util.SendResponse(c, http.StatusOK, assets)
}

// Update godoc
// @Summary Update asset
// @Description Update an existing asset by ID
// @Tags assets
// @Accept json
// @Produce json
// @Param assetId path int64 true "Asset ID"
// @Param payload body dto.UpdateAssetRequest true "Update asset payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/assets/{assetId} [put]
// @Security BearerAuth
func (h *AssetHandler) Update(c *gin.Context) {
	slog.Info("request started", "handler", "AssetHandler.Update", "method", c.Request.Method, "path", c.Request.URL.Path)

	id, err := parseIntegerID(c, "assetId")
	if err != nil {
		slog.Warn("invalid asset ID", "handler", "AssetHandler.Update", "assetId", c.Param("assetId"))
		util.SendErrorResponse(c, http.StatusBadRequest, "invalid asset id")
		return
	}

	var req dto.UpdateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			slog.Warn("validation failed", "handler", "AssetHandler.Update", "error", ve)
			util.SendErrorResponse(c, http.StatusBadRequest, util.FormatValidationErrors(err))
			return
		}
		slog.Warn("failed to bind request body", "handler", "AssetHandler.Update", "error", err)
		util.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Update(c.Request.Context(), id, &req); err != nil {
		util.HandleError(c, err, "AssetHandler.Update")
		return
	}

	slog.Info("asset updated", "handler", "AssetHandler.Update", "assetID", id)
	util.SendResponse(c, http.StatusOK, map[string]string{"message": "asset updated successfully"})
}

// Delete godoc
// @Summary Delete asset
// @Description Delete an asset by ID
// @Tags assets
// @Produce json
// @Param assetId path int64 true "Asset ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/assets/{assetId} [delete]
// @Security BearerAuth
func (h *AssetHandler) Delete(c *gin.Context) {
	slog.Info("request started", "handler", "AssetHandler.Delete", "method", c.Request.Method, "path", c.Request.URL.Path)

	id, err := parseIntegerID(c, "assetId")
	if err != nil {
		slog.Warn("invalid asset ID", "handler", "AssetHandler.Delete", "assetId", c.Param("assetId"))
		util.SendErrorResponse(c, http.StatusBadRequest, "invalid asset id")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		util.HandleError(c, err, "AssetHandler.Delete")
		return
	}

	slog.Info("asset deleted", "handler", "AssetHandler.Delete", "assetID", id)
	util.SendResponse(c, http.StatusOK, map[string]string{"message": "asset deleted successfully"})
}

func (h *AssetHandler) SetRoutes(r *gin.RouterGroup) {
	assets := r.Group("/assets")
	assets.POST("", h.Create)
	assets.GET("", h.GetAll)
	assets.GET("/:assetId", h.GetByID)
	assets.PUT("/:assetId", h.Update)
	assets.DELETE("/:assetId", h.Delete)
}
