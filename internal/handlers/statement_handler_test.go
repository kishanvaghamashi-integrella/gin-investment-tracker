package handler_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	handler "gin-investment-tracker/internal/handlers"
	middleware "gin-investment-tracker/internal/middlewares"
	"gin-investment-tracker/internal/mocks"
	"gin-investment-tracker/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupStatementRouter(svc *mocks.MockStatementService) *gin.Engine {
	r := gin.New()
	h := handler.NewCasStatementHandler(svc)
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth())
	h.SetRoutes(protected)
	return r
}

func validStatementToken(t *testing.T, userID int64, email string) string {
	t.Helper()
	t.Setenv("JWT_SECRET", "test-secret-key")
	token, err := util.GenerateToken(userID, email)
	require.NoError(t, err)
	return token
}

// multipartBody creates a multipart form body with the given file content.
func multipartBody(t *testing.T, filename string, content []byte, extraFields map[string]string) (*bytes.Buffer, string) {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	for k, v := range extraFields {
		err := writer.WriteField(k, v)
		require.NoError(t, err)
	}

	part, err := writer.CreateFormFile("file", filename)
	require.NoError(t, err)
	_, err = part.Write(content)
	require.NoError(t, err)

	require.NoError(t, writer.Close())
	return &body, writer.FormDataContentType()
}

// ─────────────────────────────────────────────
// POST /api/cas-statement — ProcessCasStatement
// ─────────────────────────────────────────────

func TestStatementHandler_ProcessCasStatement_NoAuth(t *testing.T) {
	svc := new(mocks.MockStatementService)
	r := setupStatementRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/cas-statement", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "ProcessCasFile")
}

func TestStatementHandler_ProcessCasStatement_InvalidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	svc := new(mocks.MockStatementService)
	r := setupStatementRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/cas-statement", nil)
	req.AddCookie(&http.Cookie{Name: "jwt_token", Value: "not.a.valid.token"})
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "ProcessCasFile")
}

func TestStatementHandler_ProcessCasStatement_NoFile(t *testing.T) {
	token := validStatementToken(t, 1, "user@example.com")
	svc := new(mocks.MockStatementService)
	r := setupStatementRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/cas-statement", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "jwt_token", Value: token})
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "No cas statement provided")
	svc.AssertNotCalled(t, "ProcessCasFile")
}

func TestStatementHandler_ProcessCasStatement_Success(t *testing.T) {
	token := validStatementToken(t, 1, "user@example.com")
	svc := new(mocks.MockStatementService)
	svc.On("ProcessCasFile", mock.Anything, mock.Anything, "", int64(1)).Return()

	r := setupStatementRouter(svc)

	body, contentType := multipartBody(t, "statement.pdf", []byte("fake pdf content"), nil)
	req := httptest.NewRequest(http.MethodPost, "/api/cas-statement", body)
	req.Header.Set("Content-Type", contentType)
	req.AddCookie(&http.Cookie{Name: "jwt_token", Value: token})
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	assert.Contains(t, w.Body.String(), "being processed")
	svc.AssertExpectations(t)
}

func TestStatementHandler_ProcessCasStatement_WithPassword(t *testing.T) {
	token := validStatementToken(t, 1, "user@example.com")
	svc := new(mocks.MockStatementService)
	svc.On("ProcessCasFile", mock.Anything, mock.Anything, "mypassword", int64(1)).Return()

	r := setupStatementRouter(svc)

	body, contentType := multipartBody(t, "protected.pdf", []byte("fake pdf content"), map[string]string{"password": "mypassword"})
	req := httptest.NewRequest(http.MethodPost, "/api/cas-statement", body)
	req.Header.Set("Content-Type", contentType)
	req.AddCookie(&http.Cookie{Name: "jwt_token", Value: token})
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	svc.AssertExpectations(t)
}
