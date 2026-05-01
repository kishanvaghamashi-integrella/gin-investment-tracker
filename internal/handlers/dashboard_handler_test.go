package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	dto "gin-investment-tracker/internal/dtos"
	handler "gin-investment-tracker/internal/handlers"
	middleware "gin-investment-tracker/internal/middlewares"
	"gin-investment-tracker/internal/mocks"
	"gin-investment-tracker/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupDashboardRouter(svc *mocks.MockDashboardService) *gin.Engine {
	r := gin.New()
	h := handler.NewDashboardHandler(svc)
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth())
	h.SetRoutes(protected)
	return r
}

func validDashboardToken(t *testing.T, userID int64, email string) string {
	t.Helper()
	t.Setenv("JWT_SECRET", "test-secret-key")
	token, err := util.GenerateToken(userID, email)
	require.NoError(t, err)
	return token
}

// ─────────────────────────────────────────────
// GET /api/dashboard — GetDashboardData
// ─────────────────────────────────────────────

func TestDashboardHandler_GetDashboardData_NoAuthHeader(t *testing.T) {
	svc := new(mocks.MockDashboardService)
	r := setupDashboardRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "GetDashboardData")
}

func TestDashboardHandler_GetDashboardData_NonBearerFormat(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockDashboardService)
	r := setupDashboardRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)
	req.Header.Set("Authorization", "Token sometoken")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "GetDashboardData")
}

func TestDashboardHandler_GetDashboardData_InvalidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockDashboardService)
	r := setupDashboardRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)
	req.Header.Set("Authorization", "Bearer not.a.valid.token")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "GetDashboardData")
}

func TestDashboardHandler_GetDashboardData_ServiceInternalError(t *testing.T) {
	token := validDashboardToken(t, 1, "user@example.com")
	svc := new(mocks.MockDashboardService)
	svc.On("GetDashboardData", mock.Anything, int64(1)).
		Return(nil, util.NewInternalError("db failure"))

	r := setupDashboardRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestDashboardHandler_GetDashboardData_Success(t *testing.T) {
	token := validDashboardToken(t, 1, "user@example.com")
	svc := new(mocks.MockDashboardService)

	desc := "initial buy"
	curr := 50000.00
	prev := 48000.00
	total := 45000.00
	svc.On("GetDashboardData", mock.Anything, int64(1)).
		Return(&dto.DashboardDataDto{
			CurrentInvestmentValue:     &curr,
			PreviousDayInvestmentValue: &prev,
			TotalInvestedValue:         &total,
			TopHoldings: []dto.HoldingResponseDto{
				{
					AssetName:           "Infosys",
					AssetInstrumentType: "stock",
					Quantity:            10,
					AveragePrice:        1500,
					CurrentPrice:        1600,
					InvestedCapital:     15000,
					CurrentCapital:      16000,
					ReturnPercentage:    6.67,
				},
			},
			RecentTransactions: []dto.TransactionResponseDto{
				{
					AssetName:           "Infosys",
					AssetInstrumentType: "stock",
					TxnType:             "buy",
					Quantity:            10,
					Price:               1500,
					Description:         &desc,
					TxnDate:             time.Now(),
				},
			},
		}, nil)

	r := setupDashboardRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "current_investment_value")
	assert.Contains(t, w.Body.String(), "Infosys")
	svc.AssertExpectations(t)
}

func TestDashboardHandler_GetDashboardData_Success_EmptyData(t *testing.T) {
	token := validDashboardToken(t, 2, "newuser@example.com")
	svc := new(mocks.MockDashboardService)
	svc.On("GetDashboardData", mock.Anything, int64(2)).
		Return(&dto.DashboardDataDto{}, nil)

	r := setupDashboardRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}
