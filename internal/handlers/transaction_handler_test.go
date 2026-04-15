package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	handler "gin-investment-tracker/internal/handlers"
	middleware "gin-investment-tracker/internal/middlewares"
	"gin-investment-tracker/internal/mocks"
	model "gin-investment-tracker/internal/models"
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

func setupTransactionRouter(svc *mocks.MockTransactionService) *gin.Engine {
	r := gin.New()
	h := handler.NewTransactionHandler(svc)
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth())
	h.SetRoutes(protected)
	return r
}

func validTxnToken(t *testing.T, userID int64, email string) string {
	t.Helper()
	t.Setenv("JWT_SECRET", "test-secret-key")
	token, err := util.GenerateToken(userID, email)
	require.NoError(t, err)
	return token
}

func txnJSONBody(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return bytes.NewBuffer(b)
}

func validCreateTxnBody() map[string]any {
	return map[string]any{
		"asset_id": 1,
		"txn_type": "BUY",
		"quantity": 10.0,
		"price":    100.0,
		"txn_date": time.Now().Format(time.RFC3339),
	}
}

// ─────────────────────────────────────────────
// POST /api/transactions — Create
// ─────────────────────────────────────────────

func TestTransactionHandler_Create_Success(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)
	svc.On("Create", mock.Anything, mock.Anything, int64(1)).
		Return(&model.Transaction{ID: 42, TxnType: "BUY", Quantity: 10, Price: 100}, nil)

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions", txnJSONBody(t, validCreateTxnBody()))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "transaction created with id 42")
	svc.AssertExpectations(t)
}

func TestTransactionHandler_Create_MalformedJSON(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewBufferString(`{invalid`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestTransactionHandler_Create_MissingAssetID(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)

	body := validCreateTxnBody()
	delete(body, "asset_id")

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions", txnJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestTransactionHandler_Create_MissingTxnType(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)

	body := validCreateTxnBody()
	delete(body, "txn_type")

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions", txnJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestTransactionHandler_Create_InvalidTxnType(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)

	body := validCreateTxnBody()
	body["txn_type"] = "HOLD"

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions", txnJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid transaction type")
	svc.AssertNotCalled(t, "Create")
}

func TestTransactionHandler_Create_MissingQuantity(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)

	body := validCreateTxnBody()
	delete(body, "quantity")

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions", txnJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestTransactionHandler_Create_MissingPrice(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)

	body := validCreateTxnBody()
	delete(body, "price")

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions", txnJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestTransactionHandler_Create_MissingTxnDate(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)

	body := validCreateTxnBody()
	delete(body, "txn_date")

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions", txnJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestTransactionHandler_Create_ServiceBadRequest(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)
	svc.On("Create", mock.Anything, mock.Anything, int64(1)).
		Return(nil, util.NewBadRequestError("Cannot sell asset that is not currently held"))

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions", txnJSONBody(t, validCreateTxnBody()))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Cannot sell asset that is not currently held")
	svc.AssertExpectations(t)
}

func TestTransactionHandler_Create_ServiceNotFound(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)
	svc.On("Create", mock.Anything, mock.Anything, int64(1)).
		Return(nil, util.NewNotFoundError("asset not found on database"))

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions", txnJSONBody(t, validCreateTxnBody()))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "asset not found on database")
	svc.AssertExpectations(t)
}

func TestTransactionHandler_Create_ServiceInternalError(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)
	svc.On("Create", mock.Anything, mock.Anything, int64(1)).
		Return(nil, util.NewInternalError("db failure"))

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions", txnJSONBody(t, validCreateTxnBody()))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestTransactionHandler_Create_NoAuthHeader(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockTransactionService)

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions", txnJSONBody(t, validCreateTxnBody()))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestTransactionHandler_Create_NonBearerFormat(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockTransactionService)

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions", txnJSONBody(t, validCreateTxnBody()))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token sometoken")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestTransactionHandler_Create_InvalidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockTransactionService)

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions", txnJSONBody(t, validCreateTxnBody()))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer not.a.valid.token")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "Create")
}

// ─────────────────────────────────────────────
// GET /api/transactions — GetAllByUserID
// ─────────────────────────────────────────────

func TestTransactionHandler_GetAllByUserID_Success(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)
	svc.On("GetAllByUserID", mock.Anything, int64(1), 50, 0).
		Return([]dto.ResponseTransactionDto{{ID: 1, AssetName: "INFY", TxnType: "BUY"}}, nil)

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "INFY")
	svc.AssertExpectations(t)
}

func TestTransactionHandler_GetAllByUserID_CustomPagination(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)
	svc.On("GetAllByUserID", mock.Anything, int64(1), 10, 20).
		Return([]dto.ResponseTransactionDto{}, nil)

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/transactions?limit=10&offset=20", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestTransactionHandler_GetAllByUserID_InvalidLimit(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/transactions?limit=abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "limit must be a positive integer")
	svc.AssertNotCalled(t, "GetAllByUserID")
}

func TestTransactionHandler_GetAllByUserID_InvalidOffset(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/transactions?offset=abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "offset must be a non-negative integer")
	svc.AssertNotCalled(t, "GetAllByUserID")
}

func TestTransactionHandler_GetAllByUserID_LimitCappedToMax(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)
	svc.On("GetAllByUserID", mock.Anything, int64(1), 200, 0).
		Return([]dto.ResponseTransactionDto{}, nil)

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/transactions?limit=999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestTransactionHandler_GetAllByUserID_ServiceNotFound(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)
	svc.On("GetAllByUserID", mock.Anything, int64(1), 50, 0).
		Return(nil, util.NewNotFoundError("user with id 1 not found"))

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	svc.AssertExpectations(t)
}

func TestTransactionHandler_GetAllByUserID_ServiceInternalError(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)
	svc.On("GetAllByUserID", mock.Anything, int64(1), 50, 0).
		Return(nil, util.NewInternalError("db failure"))

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestTransactionHandler_GetAllByUserID_NoAuthHeader(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockTransactionService)

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "GetAllByUserID")
}

// ─────────────────────────────────────────────
// PUT /api/transactions/:txnId — Update
// ─────────────────────────────────────────────

func TestTransactionHandler_Update_Success(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)
	svc.On("Update", mock.Anything, int64(1), mock.Anything).Return(nil)

	r := setupTransactionRouter(svc)
	body := txnJSONBody(t, map[string]any{"quantity": 20.0})
	req := httptest.NewRequest(http.MethodPut, "/api/transactions/1", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "transaction updated successfully")
	svc.AssertExpectations(t)
}

func TestTransactionHandler_Update_NonNumericID(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)

	r := setupTransactionRouter(svc)
	body := txnJSONBody(t, map[string]any{"quantity": 20.0})
	req := httptest.NewRequest(http.MethodPut, "/api/transactions/abc", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid transaction id")
	svc.AssertNotCalled(t, "Update")
}

func TestTransactionHandler_Update_MalformedJSON(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodPut, "/api/transactions/1", bytes.NewBufferString(`{invalid`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Update")
}

func TestTransactionHandler_Update_InvalidTxnType(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)

	r := setupTransactionRouter(svc)
	body := txnJSONBody(t, map[string]any{"txn_type": "HOLD"})
	req := httptest.NewRequest(http.MethodPut, "/api/transactions/1", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid transaction type")
	svc.AssertNotCalled(t, "Update")
}

func TestTransactionHandler_Update_NotFound(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)
	svc.On("Update", mock.Anything, int64(99), mock.Anything).
		Return(util.NewNotFoundError("transaction with id 99 not found"))

	r := setupTransactionRouter(svc)
	body := txnJSONBody(t, map[string]any{"quantity": 5.0})
	req := httptest.NewRequest(http.MethodPut, "/api/transactions/99", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "transaction with id 99 not found")
	svc.AssertExpectations(t)
}

func TestTransactionHandler_Update_InternalError(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)
	svc.On("Update", mock.Anything, int64(1), mock.Anything).
		Return(util.NewInternalError("db failure"))

	r := setupTransactionRouter(svc)
	body := txnJSONBody(t, map[string]any{"quantity": 5.0})
	req := httptest.NewRequest(http.MethodPut, "/api/transactions/1", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestTransactionHandler_Update_NoAuthHeader(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockTransactionService)

	r := setupTransactionRouter(svc)
	body := txnJSONBody(t, map[string]any{"quantity": 5.0})
	req := httptest.NewRequest(http.MethodPut, "/api/transactions/1", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "Update")
}

// ─────────────────────────────────────────────
// DELETE /api/transactions/:txnId — Delete
// ─────────────────────────────────────────────

func TestTransactionHandler_Delete_Success(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)
	svc.On("Delete", mock.Anything, int64(1)).Return(nil)

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/transactions/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "transaction deleted successfully")
	svc.AssertExpectations(t)
}

func TestTransactionHandler_Delete_NonNumericID(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/transactions/abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid transaction id")
	svc.AssertNotCalled(t, "Delete")
}

func TestTransactionHandler_Delete_NotFound(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)
	svc.On("Delete", mock.Anything, int64(99)).
		Return(util.NewNotFoundError("transaction with id 99 not found"))

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/transactions/99", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "transaction with id 99 not found")
	svc.AssertExpectations(t)
}

func TestTransactionHandler_Delete_InternalError(t *testing.T) {
	token := validTxnToken(t, 1, "user@example.com")
	svc := new(mocks.MockTransactionService)
	svc.On("Delete", mock.Anything, int64(1)).
		Return(util.NewInternalError("db failure"))

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/transactions/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestTransactionHandler_Delete_NoAuthHeader(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockTransactionService)

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/transactions/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "Delete")
}

func TestTransactionHandler_Delete_NonBearerFormat(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockTransactionService)

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/transactions/1", nil)
	req.Header.Set("Authorization", "Token sometoken")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "Delete")
}

func TestTransactionHandler_Delete_InvalidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockTransactionService)

	r := setupTransactionRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/transactions/1", nil)
	req.Header.Set("Authorization", "Bearer not.a.valid.token")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "Delete")
}
