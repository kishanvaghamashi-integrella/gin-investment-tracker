package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	handler "gin-investment-tracker/internal/handlers"
	middleware "gin-investment-tracker/internal/middlewares"
	"gin-investment-tracker/internal/mocks"
	model "gin-investment-tracker/internal/models"
	"gin-investment-tracker/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// setupRouter wires up a handler with the given mock service and registers
// all four routes, mirroring routes/routes.go.
func setupRouter(svc *mocks.MockUserService) *gin.Engine {
	r := gin.New()
	h := handler.NewUserHandler(svc)
	api := r.Group("/api")
	users := api.Group("/users")
	users.POST("", h.Create)
	users.POST("/login", h.Login)
	users.GET("/verify", middleware.JWTAuth(), h.Verify)
	users.DELETE("", middleware.JWTAuth(), h.Delete)
	return r
}

// jsonBody encodes v to JSON and returns a *bytes.Buffer.
func jsonBody(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return bytes.NewBuffer(b)
}

// validToken generates a real JWT for testing protected endpoints.
func validToken(t *testing.T, userID int64, email string) string {
	t.Helper()
	t.Setenv("JWT_SECRET", "test-secret-key")
	token, err := util.GenerateToken(userID, email)
	require.NoError(t, err)
	return token
}

// ─────────────────────────────────────────────
// POST /api/users — Create
// ─────────────────────────────────────────────

func TestUserHandler_Create_Success(t *testing.T) {
	svc := new(mocks.MockUserService)
	svc.On("Create", mock.Anything, mock.Anything).Return(nil)

	r := setupRouter(svc)
	body := jsonBody(t, map[string]string{
		"name":     "Alice",
		"email":    "alice@example.com",
		"password": "secret123",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "user created successfully")
	svc.AssertExpectations(t)
}

func TestUserHandler_Create_MalformedJSON(t *testing.T) {
	svc := new(mocks.MockUserService)
	r := setupRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBufferString(`{invalid`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestUserHandler_Create_MissingName(t *testing.T) {
	svc := new(mocks.MockUserService)
	r := setupRouter(svc)

	body := jsonBody(t, map[string]string{
		"email":    "alice@example.com",
		"password": "secret123",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestUserHandler_Create_NameTooShort(t *testing.T) {
	svc := new(mocks.MockUserService)
	r := setupRouter(svc)

	body := jsonBody(t, map[string]string{
		"name":     "Al",
		"email":    "alice@example.com",
		"password": "secret123",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestUserHandler_Create_MissingEmail(t *testing.T) {
	svc := new(mocks.MockUserService)
	r := setupRouter(svc)

	body := jsonBody(t, map[string]string{
		"name":     "Alice",
		"password": "secret123",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestUserHandler_Create_InvalidEmailFormat(t *testing.T) {
	svc := new(mocks.MockUserService)
	r := setupRouter(svc)

	body := jsonBody(t, map[string]string{
		"name":     "Alice",
		"email":    "not-an-email",
		"password": "secret123",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestUserHandler_Create_MissingPassword(t *testing.T) {
	svc := new(mocks.MockUserService)
	r := setupRouter(svc)

	body := jsonBody(t, map[string]string{
		"name":  "Alice",
		"email": "alice@example.com",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestUserHandler_Create_PasswordTooShort(t *testing.T) {
	svc := new(mocks.MockUserService)
	r := setupRouter(svc)

	body := jsonBody(t, map[string]string{
		"name":     "Alice",
		"email":    "alice@example.com",
		"password": "abc",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestUserHandler_Create_ServiceBadRequest(t *testing.T) {
	svc := new(mocks.MockUserService)
	svc.On("Create", mock.Anything, mock.Anything).Return(util.NewBadRequestError("email already exists"))

	r := setupRouter(svc)
	body := jsonBody(t, map[string]string{
		"name":     "Alice",
		"email":    "alice@example.com",
		"password": "secret123",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "email already exists")
	svc.AssertExpectations(t)
}

func TestUserHandler_Create_ServiceInternalError(t *testing.T) {
	svc := new(mocks.MockUserService)
	svc.On("Create", mock.Anything, mock.Anything).Return(util.NewInternalError("db failure"))

	r := setupRouter(svc)
	body := jsonBody(t, map[string]string{
		"name":     "Alice",
		"email":    "alice@example.com",
		"password": "secret123",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

// ─────────────────────────────────────────────
// POST /api/users/login — Login
// ─────────────────────────────────────────────

func TestUserHandler_Login_Success(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockUserService)
	svc.On("Login", mock.Anything, mock.Anything).Return(&model.User{
		ID:    1,
		Name:  "Alice",
		Email: "alice@example.com",
	}, nil)

	r := setupRouter(svc)
	body := jsonBody(t, map[string]string{
		"email":    "alice@example.com",
		"password": "secret123",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "login successful", resp["message"])
	userField, ok := resp["user"].(map[string]any)
	require.True(t, ok, "user field must be an object")
	assert.NotEmpty(t, userField["token"], "token must be present in the response")
	svc.AssertExpectations(t)
}

func TestUserHandler_Login_MissingEmail(t *testing.T) {
	svc := new(mocks.MockUserService)
	r := setupRouter(svc)

	body := jsonBody(t, map[string]string{"password": "secret123"})
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Login")
}

func TestUserHandler_Login_InvalidEmailFormat(t *testing.T) {
	svc := new(mocks.MockUserService)
	r := setupRouter(svc)

	body := jsonBody(t, map[string]string{"email": "not-an-email", "password": "secret123"})
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Login")
}

func TestUserHandler_Login_PasswordTooShort(t *testing.T) {
	svc := new(mocks.MockUserService)
	r := setupRouter(svc)

	body := jsonBody(t, map[string]string{"email": "alice@example.com", "password": "abc"})
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Login")
}

func TestUserHandler_Login_ServiceBadRequest(t *testing.T) {
	svc := new(mocks.MockUserService)
	svc.On("Login", mock.Anything, mock.Anything).Return(nil, util.NewBadRequestError("invalid email or password"))

	r := setupRouter(svc)
	body := jsonBody(t, map[string]string{"email": "alice@example.com", "password": "wrongpass"})
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid email or password")
	svc.AssertExpectations(t)
}

func TestUserHandler_Login_ServiceInternalError(t *testing.T) {
	svc := new(mocks.MockUserService)
	svc.On("Login", mock.Anything, mock.Anything).Return(nil, util.NewInternalError("db failure"))

	r := setupRouter(svc)
	body := jsonBody(t, map[string]string{"email": "alice@example.com", "password": "secret123"})
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

// ─────────────────────────────────────────────
// GET /api/users/verify — Verify (JWT protected)
// ─────────────────────────────────────────────

func TestUserHandler_Verify_Success(t *testing.T) {
	token := validToken(t, 1, "alice@example.com")
	svc := new(mocks.MockUserService)
	svc.On("GetByID", mock.Anything, int64(1)).Return(&model.User{
		ID:    1,
		Name:  "Alice",
		Email: "alice@example.com",
	}, nil)

	r := setupRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/users/verify", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "token is valid", resp["message"])
	userField, ok := resp["user"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "alice@example.com", userField["email"])
	svc.AssertExpectations(t)
}

func TestUserHandler_Verify_NoAuthHeader(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockUserService)
	r := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/users/verify", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "missing authorization header")
	svc.AssertNotCalled(t, "GetByID")
}

func TestUserHandler_Verify_NonBearerFormat(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockUserService)
	r := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/users/verify", nil)
	req.Header.Set("Authorization", "Token some-random-token")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid authorization header format")
	svc.AssertNotCalled(t, "GetByID")
}

func TestUserHandler_Verify_InvalidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockUserService)
	r := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/users/verify", nil)
	req.Header.Set("Authorization", "Bearer this.is.not.a.valid.jwt")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "GetByID")
}

func TestUserHandler_Verify_UserNotFound(t *testing.T) {
	token := validToken(t, 99, "ghost@example.com")
	svc := new(mocks.MockUserService)
	svc.On("GetByID", mock.Anything, int64(99)).Return(nil, util.NewNotFoundError("User not found"))

	r := setupRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/users/verify", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "User not found")
	svc.AssertExpectations(t)
}

func TestUserHandler_Verify_ServiceInternalError(t *testing.T) {
	token := validToken(t, 1, "alice@example.com")
	svc := new(mocks.MockUserService)
	svc.On("GetByID", mock.Anything, int64(1)).Return(nil, util.NewInternalError("db failure"))

	r := setupRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/users/verify", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

// ─────────────────────────────────────────────
// DELETE /api/users — Delete (JWT protected)
// ─────────────────────────────────────────────

func TestUserHandler_Delete_Success(t *testing.T) {
	token := validToken(t, 1, "alice@example.com")
	svc := new(mocks.MockUserService)
	svc.On("Delete", mock.Anything, int64(1)).Return(nil)

	r := setupRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "user deleted successfully")
	svc.AssertExpectations(t)
}

func TestUserHandler_Delete_NoAuthHeader(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockUserService)
	r := setupRouter(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/users", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "missing authorization header")
	svc.AssertNotCalled(t, "Delete")
}

func TestUserHandler_Delete_InvalidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockUserService)
	r := setupRouter(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/users", nil)
	req.Header.Set("Authorization", "Bearer completely.wrong.token")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "Delete")
}

func TestUserHandler_Delete_UserNotFound(t *testing.T) {
	token := validToken(t, 99, "ghost@example.com")
	svc := new(mocks.MockUserService)
	svc.On("Delete", mock.Anything, int64(99)).Return(util.NewNotFoundError("no user found with id 99"))

	r := setupRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "no user found with id 99")
	svc.AssertExpectations(t)
}

func TestUserHandler_Delete_ServiceInternalError(t *testing.T) {
	token := validToken(t, 1, "alice@example.com")
	svc := new(mocks.MockUserService)
	svc.On("Delete", mock.Anything, int64(1)).Return(util.NewInternalError("db failure"))

	r := setupRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}
