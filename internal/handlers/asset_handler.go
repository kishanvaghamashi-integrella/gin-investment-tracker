package handler

import (
	"errors"
	"fmt"
	dto "gin-investment-tracker/internal/dtos"
	service "gin-investment-tracker/internal/services"
	"gin-investment-tracker/internal/util"
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

func (h *AssetHandler) SetRoutes(r *gin.RouterGroup) {
	assets := r.Group("/assets")
	assets.POST("", h.Create)
	assets.GET("", h.GetAll)
	assets.GET("/:assetId", h.GetByID)
	assets.PUT("/:assetId", h.Update)
	assets.DELETE("/:assetId", h.Delete)
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
// @Security CookieAuth
func (h *AssetHandler) Create(c *gin.Context) {
	log := util.FromContext(c.Request.Context()).With("handler", "AssetHandler.Create")
	log.Infow("request started", "method", c.Request.Method, "path", c.Request.URL.Path)

	var req dto.CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			log.Warnw("validation failed", "error", ve)
			util.SendErrorResponse(c, http.StatusBadRequest, util.FormatValidationErrors(err))
			return
		}
		log.Warnw("failed to bind request body", "error", err)
		util.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	asset, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		util.HandleError(c, err, log)
		return
	}

	log.Infow("asset created", "asset_id", asset.ID)
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
// @Security CookieAuth
func (h *AssetHandler) GetByID(c *gin.Context) {
	log := util.FromContext(c.Request.Context()).With("handler", "AssetHandler.GetByID")
	log.Infow("request started", "method", c.Request.Method, "path", c.Request.URL.Path)

	id, err := parseIntegerID(c, "assetId")
	if err != nil {
		log.Warnw("invalid asset ID", "asset_id", c.Param("assetId"))
		util.SendErrorResponse(c, http.StatusBadRequest, "invalid asset id")
		return
	}

	asset, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		util.HandleError(c, err, log)
		return
	}

	log.Infow("asset retrieved", "asset_id", id)
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
// @Security CookieAuth
func (h *AssetHandler) GetAll(c *gin.Context) {
	log := util.FromContext(c.Request.Context()).With("handler", "AssetHandler.GetAll")
	log.Infow("request started", "method", c.Request.Method, "path", c.Request.URL.Path)

	limit, offset, err := parsePaginationParams(c)
	if err != nil {
		log.Warnw("invalid pagination params", "error", err)
		util.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	assets, err := h.service.GetAll(c.Request.Context(), limit, offset)
	if err != nil {
		util.HandleError(c, err, log)
		return
	}

	log.Infow("assets retrieved", "count", len(assets))
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
// @Security CookieAuth
func (h *AssetHandler) Update(c *gin.Context) {
	log := util.FromContext(c.Request.Context()).With("handler", "AssetHandler.Update")
	log.Infow("request started", "method", c.Request.Method, "path", c.Request.URL.Path)

	id, err := parseIntegerID(c, "assetId")
	if err != nil {
		log.Warnw("invalid asset ID", "asset_id", c.Param("assetId"))
		util.SendErrorResponse(c, http.StatusBadRequest, "invalid asset id")
		return
	}

	var req dto.UpdateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			log.Warnw("validation failed", "error", ve)
			util.SendErrorResponse(c, http.StatusBadRequest, util.FormatValidationErrors(err))
			return
		}
		log.Warnw("failed to bind request body", "error", err)
		util.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Update(c.Request.Context(), id, &req); err != nil {
		util.HandleError(c, err, log)
		return
	}

	log.Infow("asset updated", "asset_id", id)
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
// @Security CookieAuth
func (h *AssetHandler) Delete(c *gin.Context) {
	log := util.FromContext(c.Request.Context()).With("handler", "AssetHandler.Delete")
	log.Infow("request started", "method", c.Request.Method, "path", c.Request.URL.Path)

	id, err := parseIntegerID(c, "assetId")
	if err != nil {
		log.Warnw("invalid asset ID", "asset_id", c.Param("assetId"))
		util.SendErrorResponse(c, http.StatusBadRequest, "invalid asset id")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		util.HandleError(c, err, log)
		return
	}

	log.Infow("asset deleted", "asset_id", id)
	util.SendResponse(c, http.StatusOK, map[string]string{"message": "asset deleted successfully"})
}
