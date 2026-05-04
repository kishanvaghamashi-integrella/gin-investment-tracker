package handler

import (
	"encoding/json"
	"errors"
	"io"

	dto "gin-investment-tracker/internal/dtos"
	middleware "gin-investment-tracker/internal/middlewares"
	service "gin-investment-tracker/internal/services"
	"gin-investment-tracker/internal/util"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthHandler struct {
	googleOAuthCfg *oauth2.Config
	service        service.AuthServiceInterface
}

func NewAuthHandler(svc service.AuthServiceInterface) *AuthHandler {
	googleOAuthConfig := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return &AuthHandler{
		service:        svc,
		googleOAuthCfg: googleOAuthConfig,
	}
}

func (h *AuthHandler) SetRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("", h.Signup)
		auth.POST("/email/login", h.Login)
		auth.GET("/google/login", h.GoogleLogin)
		auth.GET("/google/callback", h.GoogleCallback)
		auth.POST("/logout", h.Logout)
		auth.Use(middleware.JWTAuth()).GET("/verify", h.GetUserDetails)
		auth.Use(middleware.JWTAuth()).DELETE("", h.DeleteUser)
	}
}

// Create godoc
// @Summary Create user
// @Description Create a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body dto.CreateUserRequest true "Create user payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/auth [post]
func (h *AuthHandler) Signup(c *gin.Context) {
	slog.Info("request started", "handler", "AuthHandler.Signup", "method", c.Request.Method, "path", c.Request.URL.Path)

	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			slog.Warn("validation failed", "error", ve)
			util.SendErrorResponse(c, http.StatusBadRequest, util.FormatValidationErrors(err))
			return
		}

		slog.Warn("failed to bind request body", "handler", "AuthHandler.Signup", "error", err)
		util.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Create(c.Request.Context(), &req); err != nil {
		util.HandleError(c, err, "AuthHandler.Signup")
		return
	}

	slog.Info("user created successfully", "handler", "AuthHandler.Signup")
	util.SendResponse(c, http.StatusOK, map[string]string{"message": "user created successfully"})
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password, returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body dto.LoginRequest true "Login payload"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/auth/email/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	slog.Info("request started", "handler", "AuthHandler.Login", "method", c.Request.Method, "path", c.Request.URL.Path)

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			slog.Warn("validation failed", "error", ve)
			util.SendErrorResponse(c, http.StatusBadRequest, util.FormatValidationErrors(err))
			return
		}

		slog.Warn("failed to bind request body", "handler", "AuthHandler.Login", "error", err)
		util.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	loginResp, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		util.HandleError(c, err, "AuthHandler.Login")
		return
	}

	setJwtCookie(c, loginResp.Token)

	slog.Info("user logged in successfully", "handler", "AuthHandler.Login", "userID", loginResp.ID)
	util.SendResponse(c, http.StatusOK, map[string]any{
		"message": "login successful",
		"user":    loginResp,
	})
}

// Verify godoc
// @Summary Verify token
// @Description Verify JWT cookie and return user info
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/auth/verify [get]
// @Security CookieAuth
func (h *AuthHandler) GetUserDetails(c *gin.Context) {
	slog.Info("request started", "handler", "AuthHandler.Verify", "method", c.Request.Method, "path", c.Request.URL.Path)

	userID, ok := util.GetUserIDFromContext(c)
	if !ok {
		slog.Warn("failed to parse user ID from context", "handler", "AuthHandler.Verify")
		util.SendErrorResponse(c, http.StatusBadRequest, "error while parsing the userId")
		return
	}

	user, err := h.service.GetByID(c.Request.Context(), userID)
	if err != nil {
		util.HandleError(c, err, "AuthHandler.Verify")
		return
	}

	slog.Info("token verified successfully", "handler", "AuthHandler.Verify", "userID", userID)
	util.SendResponse(c, http.StatusOK, map[string]any{
		"message": "token is valid",
		"user": dto.LoginResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	})
}

// Delete godoc
// @Summary Delete user
// @Description Delete user by ID
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorBody
// @Failure 404 {object} util.ErrorBody
// @Failure 500 {object} util.ErrorBody
// @Router /api/auth [delete]
// @Security CookieAuth
func (h *AuthHandler) DeleteUser(c *gin.Context) {
	slog.Info("request started", "handler", "AuthHandler.Delete", "method", c.Request.Method, "path", c.Request.URL.Path)

	userID, ok := util.GetUserIDFromContext(c)
	if !ok {
		slog.Warn("failed to parse user ID from context", "handler", "AuthHandler.Delete")
		util.SendErrorResponse(c, http.StatusBadRequest, "error while parsing the userId")
		return
	}

	if err := h.service.Delete(c.Request.Context(), userID); err != nil {
		util.HandleError(c, err, "AuthHandler.Delete")
		return
	}

	slog.Info("user deleted successfully", "handler", "AuthHandler.Delete", "userID", userID)
	util.SendResponse(c, http.StatusOK, map[string]string{"message": "user deleted successfully"})
}

// GoogleLogin godoc
// @Summary Initiate Google OAuth login
// @Description Redirects the user to Google's OAuth 2.0 consent screen
// @Tags auth
// @Produce json
// @Success 307 {string} string "Redirect to Google OAuth"
// @Router /api/auth/google/login [get]
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	state := uuid.New().String()
	if err := util.StoreOAuthState(state); err != nil {
		slog.Warn("oauth state store full, rejecting login initiation", "handler", "AuthHandler.GoogleLogin", "error", err)
		util.SendErrorResponse(c, http.StatusServiceUnavailable, "too many pending login attempts, please try again later")
		return
	}
	authURL := h.googleOAuthCfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	slog.Info("request started", "handler", "AuthHandler.GoogleCallback", "method", c.Request.Method, "path", c.Request.URL.Path)

	code := c.Query("code")
	if code == "" {
		slog.Warn("missing oauth code", "handler", "AuthHandler.GoogleCallback")
		c.Redirect(http.StatusSeeOther, "http://localhost:3000/login")
		return
	}

	state := c.Query("state")
	if !util.ValidateAndConsumeOAuthState(state) {
		slog.Warn("invalid or expired oauth state", "handler", "AuthHandler.GoogleCallback")
		c.Redirect(http.StatusSeeOther, "http://localhost:3000/login")
		return
	}

	token, err := h.googleOAuthCfg.Exchange(c.Request.Context(), code)
	if err != nil {
		slog.Error("failed to exchange oauth code", "handler", "AuthHandler.GoogleCallback", "error", err)
		c.Redirect(http.StatusSeeOther, "http://localhost:3000/login")
		return
	}

	client := h.googleOAuthCfg.Client(c.Request.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		slog.Error("failed to fetch user info from google", "handler", "AuthHandler.GoogleCallback", "error", err)
		c.Redirect(http.StatusSeeOther, "http://localhost:3000/login")
		return
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			slog.Error("google userinfo returned non-success status and response body could not be read", "handler", "AuthHandler.GoogleCallback", "statusCode", resp.StatusCode, "error", readErr)
			util.SendErrorResponse(c, http.StatusBadGateway, "failed to fetch user info from google")
			return
		}
		slog.Error("google userinfo returned non-success status", "handler", "AuthHandler.GoogleCallback", "statusCode", resp.StatusCode, "responseBody", string(body))
		c.Redirect(http.StatusSeeOther, "http://localhost:3000/login")
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("failed to read google userinfo response", "handler", "AuthHandler.GoogleCallback", "error", err)
		c.Redirect(http.StatusSeeOther, "http://localhost:3000/login")
		return
	}

	var userInfo dto.GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		slog.Error("failed to parse google userinfo", "handler", "AuthHandler.GoogleCallback", "error", err)
		c.Redirect(http.StatusSeeOther, "http://localhost:3000/login")
		return
	}

	if !userInfo.EmailVerified {
		slog.Warn("google email not verified", "handler", "AuthHandler.GoogleCallback", "email", userInfo.Email)
		c.Redirect(http.StatusSeeOther, "http://localhost:3000/login")
		return
	}

	loginResp, err := h.service.GoogleLogin(c.Request.Context(), &userInfo)
	if err != nil {
		util.HandleError(c, err, "AuthHandler.GoogleCallback")
		return
	}

	setJwtCookie(c, loginResp.Token)

	slog.Info("google login successful", "handler", "AuthHandler.GoogleCallback", "userID", loginResp.ID)
	c.Redirect(http.StatusSeeOther, "http://localhost:3000/dashboard")
}

// Logout godoc
// @Summary Logout
// @Description Clears the JWT auth cookie
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	slog.Info("request started", "handler", "AuthHandler.Logout", "method", c.Request.Method, "path", c.Request.URL.Path)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("jwt_token", "", -1, "/", "", false, true)
	slog.Info("user logged out successfully", "handler", "AuthHandler.Logout")
	util.SendResponse(c, http.StatusOK, map[string]string{"message": "logged out successfully"})
}

func setJwtCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("jwt_token", token, 86400, "/", "", false, true)
}
