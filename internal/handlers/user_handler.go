package handler

import (
	"errors"
	dto "gin-investment-tracker/internal/dtos"
	middleware "gin-investment-tracker/internal/middlewares"
	service "gin-investment-tracker/internal/services"
	"gin-investment-tracker/internal/util"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Create godoc
// @Summary Create user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param payload body dto.CreateUserRequest true "Create user payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/users/ [post]
func (h *UserHandler) Create(c *gin.Context) {
	slog.Info("request started", "handler", "UserHandler.Create", "method", c.Request.Method, "path", c.Request.URL.Path)

	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			slog.Warn("validation failed", "error", ve)
			util.SendErrorResponse(c, http.StatusBadRequest, util.FormatValidationErrors(err))
			return
		}

		slog.Warn("failed to bind request body", "handler", "UserHandler.Create", "error", err)
		util.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Create(c.Request.Context(), &req); err != nil {
		util.HandleError(c, err, "UserHandler.Create")
		return
	}

	slog.Info("user created successfully", "handler", "UserHandler.Create")
	util.SendResponse(c, http.StatusOK, map[string]string{"message": "user created successfully"})
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param payload body dto.LoginRequest true "Login payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	slog.Info("request started", "handler", "UserHandler.Login", "method", c.Request.Method, "path", c.Request.URL.Path)

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			slog.Warn("validation failed", "error", ve)
			util.SendErrorResponse(c, http.StatusBadRequest, util.FormatValidationErrors(err))
			return
		}

		slog.Warn("failed to bind request body", "handler", "UserHandler.Login", "error", err)
		util.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		util.HandleError(c, err, "UserHandler.Login")
		return
	}

	token, err := util.GenerateToken(user.ID, user.Email)
	if err != nil {
		slog.Error("failed to generate token", "handler", "UserHandler.Login", "userID", user.ID, "error", err)
		util.SendErrorResponse(c, http.StatusInternalServerError, "failed to generate token")
		return
	}

	slog.Info("user logged in successfully", "handler", "UserHandler.Login", "userID", user.ID)
	util.SendResponse(c, http.StatusOK, map[string]any{
		"message": "login successful",
		"user": dto.LoginResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Token: token,
		},
	})
}

// Verify godoc
// @Summary Verify token
// @Description Verify bearer token and return user info
// @Tags users
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/users/verify [get]
// @Security BearerAuth
func (h *UserHandler) Verify(c *gin.Context) {
	slog.Info("request started", "handler", "UserHandler.Verify", "method", c.Request.Method, "path", c.Request.URL.Path)

	userID, ok := util.GetUserIDFromContext(c)
	if !ok {
		slog.Warn("failed to parse user ID from context", "handler", "UserHandler.Verify")
		util.SendErrorResponse(c, http.StatusBadRequest, "error while parsing the userId")
		return
	}

	user, err := h.service.GetByID(c.Request.Context(), userID)
	if err != nil {
		util.HandleError(c, err, "UserHandler.Verify")
		return
	}

	token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

	slog.Info("token verified successfully", "handler", "UserHandler.Verify", "userID", userID)
	util.SendResponse(c, http.StatusOK, map[string]any{
		"message": "token is valid",
		"user": dto.LoginResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Token: token,
		},
	})
}

// Delete godoc
// @Summary Delete user
// @Description Delete user by ID
// @Tags users
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/users [delete]
// @Security BearerAuth
func (h *UserHandler) Delete(c *gin.Context) {
	slog.Info("request started", "handler", "UserHandler.Delete", "method", c.Request.Method, "path", c.Request.URL.Path)

	userID, ok := util.GetUserIDFromContext(c)
	if !ok {
		slog.Warn("failed to parse user ID from context", "handler", "UserHandler.Delete")
		util.SendErrorResponse(c, http.StatusBadRequest, "error while parsing the userId")
		return
	}

	if err := h.service.Delete(c.Request.Context(), userID); err != nil {
		util.HandleError(c, err, "UserHandler.Delete")
		return
	}

	slog.Info("user deleted successfully", "handler", "UserHandler.Delete", "userID", userID)
	util.SendResponse(c, http.StatusOK, map[string]string{"message": "user deleted successfully"})
}

func (h *UserHandler) SetRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.POST("", h.Create)
		users.POST("/login", h.Login)
		users.Use(middleware.JWTAuth()).GET("/verify", h.Verify)
		users.Use(middleware.JWTAuth()).DELETE("", h.Delete)
	}
}
