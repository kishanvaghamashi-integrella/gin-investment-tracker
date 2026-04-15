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

// setupAssetRouter wires up the asset handler using SetRoutes, mirroring routes/routes.go.
// JWT is applied at the group level, matching the protectedRouter pattern in routes.go.
func setupAssetRouter(svc *mocks.MockAssetService) *gin.Engine {
	r := gin.New()
	h := handler.NewAssetHandler(svc)
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth())
	h.SetRoutes(protected)
	return r
}

// validAssetToken generates a real JWT for testing protected endpoints.
func validAssetToken(t *testing.T, userID int64, email string) string {
	t.Helper()
	t.Setenv("JWT_SECRET", "test-secret-key")
	token, err := util.GenerateToken(userID, email)
	require.NoError(t, err)
	return token
}

func assetJSONBody(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return bytes.NewBuffer(b)
}

// validCreateAssetBody returns a map with all required fields set to valid values.
func validCreateAssetBody() map[string]string {
	return map[string]string{
		"symbol":          "INFY",
		"name":            "Infosys Ltd",
		"instrument_type": "stock",
		"isin":            "INE009A01021",
		"exchange":        "NSE",
	}
}

// ─────────────────────────────────────────────
// POST /api/assets — Create
// ─────────────────────────────────────────────

func TestAssetHandler_Create_Success(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)
	svc.On("Create", mock.Anything, mock.Anything).Return(&model.Asset{ID: 1, Symbol: "INFY"}, nil)

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/assets", assetJSONBody(t, validCreateAssetBody()))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "asset created with id")
	svc.AssertExpectations(t)
}

func TestAssetHandler_Create_MalformedJSON(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/assets", bytes.NewBufferString(`{invalid`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestAssetHandler_Create_MissingSymbol(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)

	body := validCreateAssetBody()
	delete(body, "symbol")

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/assets", assetJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestAssetHandler_Create_MissingName(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)

	body := validCreateAssetBody()
	delete(body, "name")

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/assets", assetJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestAssetHandler_Create_MissingInstrumentType(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)

	body := validCreateAssetBody()
	delete(body, "instrument_type")

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/assets", assetJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestAssetHandler_Create_MissingISIN(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)

	body := validCreateAssetBody()
	delete(body, "isin")

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/assets", assetJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestAssetHandler_Create_MissingExchange(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)

	body := validCreateAssetBody()
	delete(body, "exchange")

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/assets", assetJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestAssetHandler_Create_InvalidInstrumentType(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)

	body := validCreateAssetBody()
	body["instrument_type"] = "option"

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/assets", assetJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid instrument type")
	svc.AssertNotCalled(t, "Create")
}

func TestAssetHandler_Create_ServiceBadRequest(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)
	svc.On("Create", mock.Anything, mock.Anything).
		Return(nil, util.NewBadRequestError("asset already exists"))

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/assets", assetJSONBody(t, validCreateAssetBody()))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "asset already exists")
	svc.AssertExpectations(t)
}

func TestAssetHandler_Create_ServiceInternalError(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)
	svc.On("Create", mock.Anything, mock.Anything).
		Return(nil, util.NewInternalError("db failure"))

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/assets", assetJSONBody(t, validCreateAssetBody()))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestAssetHandler_Create_NoAuthHeader(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockAssetService)

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/assets", assetJSONBody(t, validCreateAssetBody()))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestAssetHandler_Create_NonBearerFormat(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockAssetService)

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/assets", assetJSONBody(t, validCreateAssetBody()))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token sometoken")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestAssetHandler_Create_InvalidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockAssetService)

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/assets", assetJSONBody(t, validCreateAssetBody()))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer not.a.valid.token")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "Create")
}

// ─────────────────────────────────────────────
// GET /api/assets/:assetId — GetByID
// ─────────────────────────────────────────────

func TestAssetHandler_GetByID_Success(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)
	svc.On("GetByID", mock.Anything, int64(42)).Return(&model.Asset{ID: 42, Symbol: "TCS"}, nil)

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/assets/42", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "TCS")
	svc.AssertExpectations(t)
}

func TestAssetHandler_GetByID_NonNumericID(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/assets/abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid asset id")
	svc.AssertNotCalled(t, "GetByID")
}

func TestAssetHandler_GetByID_NotFound(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)
	svc.On("GetByID", mock.Anything, int64(99)).
		Return(nil, util.NewNotFoundError("asset not found"))

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/assets/99", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "asset not found")
	svc.AssertExpectations(t)
}

func TestAssetHandler_GetByID_InternalError(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)
	svc.On("GetByID", mock.Anything, int64(1)).
		Return(nil, util.NewInternalError("db failure"))

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/assets/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestAssetHandler_GetByID_NoAuthHeader(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockAssetService)

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/assets/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "GetByID")
}

// ─────────────────────────────────────────────
// GET /api/assets — GetAll
// ─────────────────────────────────────────────

func TestAssetHandler_GetAll_Success(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)
	svc.On("GetAll", mock.Anything, 50, 0).Return([]model.Asset{{ID: 1, Symbol: "INFY"}}, nil)

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/assets", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "INFY")
	svc.AssertExpectations(t)
}

func TestAssetHandler_GetAll_CustomPagination(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)
	svc.On("GetAll", mock.Anything, 10, 20).Return([]model.Asset{}, nil)

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/assets?limit=10&offset=20", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestAssetHandler_GetAll_InvalidLimit(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/assets?limit=abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "limit must be a positive integer")
	svc.AssertNotCalled(t, "GetAll")
}

func TestAssetHandler_GetAll_InvalidOffset(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/assets?offset=abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "offset must be a non-negative integer")
	svc.AssertNotCalled(t, "GetAll")
}

func TestAssetHandler_GetAll_LimitCappedToMax(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)
	// helper silently caps limit=201 to 200 and calls service normally
	svc.On("GetAll", mock.Anything, 200, 0).Return([]model.Asset{}, nil)

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/assets?limit=201", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestAssetHandler_GetAll_InternalError(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)
	svc.On("GetAll", mock.Anything, 50, 0).
		Return(nil, util.NewInternalError("db failure"))

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/assets", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestAssetHandler_GetAll_NoAuthHeader(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockAssetService)

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/assets", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "GetAll")
}

// ─────────────────────────────────────────────
// PUT /api/assets/:assetId — Update
// ─────────────────────────────────────────────

func TestAssetHandler_Update_Success(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)
	svc.On("Update", mock.Anything, int64(1), mock.Anything).Return(nil)

	r := setupAssetRouter(svc)
	body := assetJSONBody(t, map[string]string{"name": "Infosys Limited"})
	req := httptest.NewRequest(http.MethodPut, "/api/assets/1", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "asset updated successfully")
	svc.AssertExpectations(t)
}

func TestAssetHandler_Update_NonNumericID(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)

	r := setupAssetRouter(svc)
	body := assetJSONBody(t, map[string]string{"name": "Infosys Limited"})
	req := httptest.NewRequest(http.MethodPut, "/api/assets/abc", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid asset id")
	svc.AssertNotCalled(t, "Update")
}

func TestAssetHandler_Update_MalformedJSON(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodPut, "/api/assets/1", bytes.NewBufferString(`{invalid`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Update")
}

func TestAssetHandler_Update_NotFound(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)
	svc.On("Update", mock.Anything, int64(99), mock.Anything).
		Return(util.NewNotFoundError("asset not found"))

	r := setupAssetRouter(svc)
	body := assetJSONBody(t, map[string]string{"name": "Infosys Limited"})
	req := httptest.NewRequest(http.MethodPut, "/api/assets/99", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "asset not found")
	svc.AssertExpectations(t)
}

func TestAssetHandler_Update_InternalError(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)
	svc.On("Update", mock.Anything, int64(1), mock.Anything).
		Return(util.NewInternalError("db failure"))

	r := setupAssetRouter(svc)
	body := assetJSONBody(t, map[string]string{"name": "Infosys Limited"})
	req := httptest.NewRequest(http.MethodPut, "/api/assets/1", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestAssetHandler_Update_NoAuthHeader(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockAssetService)

	r := setupAssetRouter(svc)
	body := assetJSONBody(t, map[string]string{"name": "Infosys Limited"})
	req := httptest.NewRequest(http.MethodPut, "/api/assets/1", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "Update")
}

// ─────────────────────────────────────────────
// DELETE /api/assets/:assetId — Delete
// ─────────────────────────────────────────────

func TestAssetHandler_Delete_Success(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)
	svc.On("Delete", mock.Anything, int64(1)).Return(nil)

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/assets/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "asset deleted successfully")
	svc.AssertExpectations(t)
}

func TestAssetHandler_Delete_NonNumericID(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/assets/abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid asset id")
	svc.AssertNotCalled(t, "Delete")
}

func TestAssetHandler_Delete_NotFound(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)
	svc.On("Delete", mock.Anything, int64(99)).
		Return(util.NewNotFoundError("asset not found"))

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/assets/99", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "asset not found")
	svc.AssertExpectations(t)
}

func TestAssetHandler_Delete_InternalError(t *testing.T) {
	token := validAssetToken(t, 1, "user@example.com")
	svc := new(mocks.MockAssetService)
	svc.On("Delete", mock.Anything, int64(1)).
		Return(util.NewInternalError("db failure"))

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/assets/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestAssetHandler_Delete_NoAuthHeader(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockAssetService)

	r := setupAssetRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/assets/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "Delete")
}
