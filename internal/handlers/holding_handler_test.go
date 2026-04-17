package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	handler "gin-investment-tracker/internal/handlers"
	middleware "gin-investment-tracker/internal/middlewares"
	"gin-investment-tracker/internal/mocks"
	"gin-investment-tracker/internal/util"

	dto "gin-investment-tracker/internal/dtos"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupHoldingRouter(svc *mocks.MockHoldingService) *gin.Engine {
	r := gin.New()
	h := handler.NewHoldingHandler(svc)
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth())
	h.SetRoutes(protected)
	return r
}

func validHoldingToken(t *testing.T, userID int64, email string) string {
	t.Helper()
	t.Setenv("JWT_SECRET", "test-secret-key")
	token, err := util.GenerateToken(userID, email)
	require.NoError(t, err)
	return token
}

// ─────────────────────────────────────────────
// GET /api/holdings — GetAll
// ─────────────────────────────────────────────

func TestHoldingHandler_GetAll_NoAuthHeader(t *testing.T) {
	svc := new(mocks.MockHoldingService)
	r := setupHoldingRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/holdings", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "GetAllByUserID")
}

func TestHoldingHandler_GetAll_NonBearerFormat(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockHoldingService)
	r := setupHoldingRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/holdings", nil)
	req.Header.Set("Authorization", "Token sometoken")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "GetAllByUserID")
}

func TestHoldingHandler_GetAll_InvalidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockHoldingService)
	r := setupHoldingRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/holdings", nil)
	req.Header.Set("Authorization", "Bearer not.a.valid.token")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "GetAllByUserID")
}

func TestHoldingHandler_GetAll_Success(t *testing.T) {
	token := validHoldingToken(t, 1, "user@example.com")
	svc := new(mocks.MockHoldingService)
	svc.On("GetAllByUserID", mock.Anything, int64(1), 50, 0).
		Return([]dto.HoldingResponseDto{
			{ID: 1, AssetName: "HDFC Flexi Cap", Quantity: 10, AveragePrice: 100},
		}, nil)

	r := setupHoldingRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/holdings", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "HDFC Flexi Cap")
	svc.AssertExpectations(t)
}

func TestHoldingHandler_GetAll_UserNotFound(t *testing.T) {
	token := validHoldingToken(t, 99, "ghost@example.com")
	svc := new(mocks.MockHoldingService)
	svc.On("GetAllByUserID", mock.Anything, int64(99), 50, 0).
		Return(nil, util.NewNotFoundError("user with id 99 not found"))

	r := setupHoldingRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/holdings", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	svc.AssertExpectations(t)
}

func TestHoldingHandler_GetAll_InternalError(t *testing.T) {
	token := validHoldingToken(t, 1, "user@example.com")
	svc := new(mocks.MockHoldingService)
	svc.On("GetAllByUserID", mock.Anything, int64(1), 50, 0).
		Return(nil, util.NewInternalError("db failure"))

	r := setupHoldingRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/holdings", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestHoldingHandler_GetAll_InvalidLimitParam(t *testing.T) {
	token := validHoldingToken(t, 1, "user@example.com")
	svc := new(mocks.MockHoldingService)

	r := setupHoldingRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/holdings?limit=abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "GetAllByUserID")
}

func TestHoldingHandler_GetAll_InvalidOffsetParam(t *testing.T) {
	token := validHoldingToken(t, 1, "user@example.com")
	svc := new(mocks.MockHoldingService)

	r := setupHoldingRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/holdings?offset=-1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "GetAllByUserID")
}

func TestHoldingHandler_GetAll_WithPaginationParams(t *testing.T) {
	token := validHoldingToken(t, 1, "user@example.com")
	svc := new(mocks.MockHoldingService)
	svc.On("GetAllByUserID", mock.Anything, int64(1), 10, 20).
		Return([]dto.HoldingResponseDto{}, nil)

	r := setupHoldingRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/holdings?limit=10&offset=20", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}
