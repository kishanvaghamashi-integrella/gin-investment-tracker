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

func setupUserAssetRouter(svc *mocks.MockUserAssetService) *gin.Engine {
	r := gin.New()
	h := handler.NewUserAssetHandler(svc)
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth())
	h.SetRoutes(protected)
	return r
}

func validUserAssetToken(t *testing.T, userID int64, email string) string {
	t.Helper()
	t.Setenv("JWT_SECRET", "test-secret-key")
	token, err := util.GenerateToken(userID, email)
	require.NoError(t, err)
	return token
}

func userAssetJSONBody(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return bytes.NewBuffer(b)
}

// ─────────────────────────────────────────────
// POST /api/user-assets — Create
// ─────────────────────────────────────────────

func TestUserAssetHandler_Create_Success(t *testing.T) {
	token := validUserAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockUserAssetService)
	svc.On("Create", mock.Anything, int64(1), mock.Anything).
		Return(&model.UserAsset{ID: 10, UserID: 1, AssetID: 5}, nil)

	r := setupUserAssetRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/user-assets", userAssetJSONBody(t, map[string]any{"asset_id": 5}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "assigned to user")
	svc.AssertExpectations(t)
}

func TestUserAssetHandler_Create_MalformedJSON(t *testing.T) {
	token := validUserAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockUserAssetService)

	r := setupUserAssetRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/user-assets", bytes.NewBufferString(`{invalid`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestUserAssetHandler_Create_MissingAssetID(t *testing.T) {
	token := validUserAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockUserAssetService)

	r := setupUserAssetRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/user-assets", userAssetJSONBody(t, map[string]any{}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestUserAssetHandler_Create_AssetIDZero(t *testing.T) {
	token := validUserAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockUserAssetService)

	r := setupUserAssetRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/user-assets", userAssetJSONBody(t, map[string]any{"asset_id": 0}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestUserAssetHandler_Create_ServiceBadRequest(t *testing.T) {
	token := validUserAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockUserAssetService)
	svc.On("Create", mock.Anything, int64(1), mock.Anything).
		Return(nil, util.NewBadRequestError("This entry already exists"))

	r := setupUserAssetRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/user-assets", userAssetJSONBody(t, map[string]any{"asset_id": 5}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "This entry already exists")
	svc.AssertExpectations(t)
}

func TestUserAssetHandler_Create_ServiceNotFound(t *testing.T) {
	token := validUserAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockUserAssetService)
	svc.On("Create", mock.Anything, int64(1), mock.Anything).
		Return(nil, util.NewNotFoundError("asset not found"))

	r := setupUserAssetRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/user-assets", userAssetJSONBody(t, map[string]any{"asset_id": 99}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	svc.AssertExpectations(t)
}

func TestUserAssetHandler_Create_ServiceInternalError(t *testing.T) {
	token := validUserAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockUserAssetService)
	svc.On("Create", mock.Anything, int64(1), mock.Anything).
		Return(nil, util.NewInternalError("db failure"))

	r := setupUserAssetRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/user-assets", userAssetJSONBody(t, map[string]any{"asset_id": 5}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestUserAssetHandler_Create_NoAuthHeader(t *testing.T) {
	svc := new(mocks.MockUserAssetService)

	r := setupUserAssetRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/user-assets", userAssetJSONBody(t, map[string]any{"asset_id": 5}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestUserAssetHandler_Create_NonBearerFormat(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockUserAssetService)

	r := setupUserAssetRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/user-assets", userAssetJSONBody(t, map[string]any{"asset_id": 5}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token sometoken")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestUserAssetHandler_Create_InvalidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockUserAssetService)

	r := setupUserAssetRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/user-assets", userAssetJSONBody(t, map[string]any{"asset_id": 5}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer invalidtoken")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "Create")
}

// ─────────────────────────────────────────────
// GET /api/user-assets — GetByUserID
// ─────────────────────────────────────────────

func TestUserAssetHandler_GetByUserID_Success(t *testing.T) {
	token := validUserAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockUserAssetService)
	svc.On("GetByUserID", mock.Anything, int64(1), 50, 0).
		Return([]model.UserAsset{{ID: 1, UserID: 1, AssetID: 5}}, nil)

	r := setupUserAssetRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/user-assets", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestUserAssetHandler_GetByUserID_NoAuthHeader(t *testing.T) {
	svc := new(mocks.MockUserAssetService)

	r := setupUserAssetRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/user-assets", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "GetByUserID")
}

func TestUserAssetHandler_GetByUserID_ServiceNotFound(t *testing.T) {
	token := validUserAssetToken(t, 99, "user@example.com")
	svc := new(mocks.MockUserAssetService)
	svc.On("GetByUserID", mock.Anything, int64(99), 50, 0).
		Return(nil, util.NewNotFoundError("user not found"))

	r := setupUserAssetRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/user-assets", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	svc.AssertExpectations(t)
}

func TestUserAssetHandler_GetByUserID_ServiceInternalError(t *testing.T) {
	token := validUserAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockUserAssetService)
	svc.On("GetByUserID", mock.Anything, int64(1), 50, 0).
		Return(nil, util.NewInternalError("db failure"))

	r := setupUserAssetRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/user-assets", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

// ─────────────────────────────────────────────
// DELETE /api/user-assets/:userAssetId — Delete
// ─────────────────────────────────────────────

func TestUserAssetHandler_Delete_Success(t *testing.T) {
	token := validUserAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockUserAssetService)
	svc.On("Delete", mock.Anything, int64(1), int64(10)).Return(nil)

	r := setupUserAssetRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/user-assets/10", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "user asset deleted successfully")
	svc.AssertExpectations(t)
}

func TestUserAssetHandler_Delete_NoAuthHeader(t *testing.T) {
	svc := new(mocks.MockUserAssetService)

	r := setupUserAssetRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/user-assets/10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "Delete")
}

func TestUserAssetHandler_Delete_InvalidUserAssetID(t *testing.T) {
	token := validUserAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockUserAssetService)

	r := setupUserAssetRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/user-assets/abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid user asset id")
	svc.AssertNotCalled(t, "Delete")
}

func TestUserAssetHandler_Delete_ServiceNotFound(t *testing.T) {
	token := validUserAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockUserAssetService)
	svc.On("Delete", mock.Anything, int64(1), int64(10)).
		Return(util.NewNotFoundError("user asset not found"))

	r := setupUserAssetRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/user-assets/10", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	svc.AssertExpectations(t)
}

func TestUserAssetHandler_Delete_ServiceInternalError(t *testing.T) {
	token := validUserAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockUserAssetService)
	svc.On("Delete", mock.Anything, int64(1), int64(10)).
		Return(util.NewInternalError("db failure"))

	r := setupUserAssetRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/user-assets/10", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}
