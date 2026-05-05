package service_test

import (
	"context"
	"errors"
	dto "gin-investment-tracker/internal/dtos"
	"gin-investment-tracker/internal/mocks"
	model "gin-investment-tracker/internal/models"
	service "gin-investment-tracker/internal/services"
	"gin-investment-tracker/internal/util"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ─────────────────────────────────────────────
// Create
// ─────────────────────────────────────────────

func TestUserService_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewAuthService(mockRepo)

	req := &dto.CreateUserRequest{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "secret123",
	}

	// Repo Create should be called once with any *model.User and return nil.
	mockRepo.On("Create", context.Background(), mock.AnythingOfType("*model.User")).Return(nil)

	err := svc.Create(context.Background(), req)

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Create_RepoError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewAuthService(mockRepo)

	req := &dto.CreateUserRequest{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "secret123",
	}

	repoErr := errors.New("duplicate key value")
	mockRepo.On("Create", context.Background(), mock.AnythingOfType("*model.User")).Return(repoErr)

	err := svc.Create(context.Background(), req)

	require.Error(t, err)
	assert.Equal(t, repoErr, err)
	mockRepo.AssertExpectations(t)
}

// ─────────────────────────────────────────────
// Login
// ─────────────────────────────────────────────

func TestUserService_Login_Success(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewAuthService(mockRepo)

	plainPassword := "secret123"
	hash, err := util.HashPassword(plainPassword)
	require.NoError(t, err)

	storedUser := &model.User{
		ID:           1,
		Name:         "Alice",
		Email:        "alice@example.com",
		PasswordHash: hash,
		IsActive:     true,
	}

	req := &dto.LoginRequest{
		Email:    "alice@example.com",
		Password: plainPassword,
	}

	mockRepo.On("GetByEmail", context.Background(), req.Email).Return(storedUser, nil)

	resp, err := svc.Login(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, storedUser.ID, resp.ID)
	assert.Equal(t, storedUser.Email, resp.Email)
	assert.NotEmpty(t, resp.Token)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_EmailNotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewAuthService(mockRepo)

	req := &dto.LoginRequest{
		Email:    "ghost@example.com",
		Password: "secret123",
	}

	mockRepo.On("GetByEmail", context.Background(), req.Email).Return(nil, pgx.ErrNoRows)

	resp, err := svc.Login(context.Background(), req)

	require.Nil(t, resp)
	require.Error(t, err)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok, "expected *util.AppError")
	assert.Equal(t, 400, appErr.Code)
	assert.Equal(t, "invalid email or password", appErr.Message)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_WrongPassword(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewAuthService(mockRepo)

	hash, err := util.HashPassword("correct-password")
	require.NoError(t, err)

	storedUser := &model.User{
		ID:           1,
		Email:        "alice@example.com",
		PasswordHash: hash,
		IsActive:     true,
	}

	req := &dto.LoginRequest{
		Email:    "alice@example.com",
		Password: "wrong-password",
	}

	mockRepo.On("GetByEmail", context.Background(), req.Email).Return(storedUser, nil)

	resp, err := svc.Login(context.Background(), req)

	require.Nil(t, resp)
	require.Error(t, err)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok, "expected *util.AppError")
	assert.Equal(t, 400, appErr.Code)
	assert.Equal(t, "invalid email or password", appErr.Message)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_RepoInternalError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewAuthService(mockRepo)

	req := &dto.LoginRequest{
		Email:    "alice@example.com",
		Password: "secret123",
	}

	mockRepo.On("GetByEmail", context.Background(), req.Email).Return(nil, errors.New("connection refused"))

	resp, err := svc.Login(context.Background(), req)

	require.Nil(t, resp)
	require.Error(t, err)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok, "expected *util.AppError")
	assert.Equal(t, 500, appErr.Code)
	assert.Equal(t, "failed to process login", appErr.Message)
	mockRepo.AssertExpectations(t)
}

// ─────────────────────────────────────────────
// GetByID
// ─────────────────────────────────────────────

func TestUserService_GetByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewAuthService(mockRepo)

	expected := &model.User{ID: 42, Name: "Bob", Email: "bob@example.com"}
	mockRepo.On("GetByID", context.Background(), int64(42)).Return(expected, nil)

	user, err := svc.GetByID(context.Background(), 42)

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, int64(42), user.ID)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewAuthService(mockRepo)

	notFound := util.NewNotFoundError("User not found")
	mockRepo.On("GetByID", context.Background(), int64(99)).Return(nil, notFound)

	user, err := svc.GetByID(context.Background(), 99)

	require.Nil(t, user)
	require.Error(t, err)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok)
	assert.Equal(t, 404, appErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetByID_InternalError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewAuthService(mockRepo)

	mockRepo.On("GetByID", context.Background(), int64(1)).Return(nil, errors.New("db timeout"))

	user, err := svc.GetByID(context.Background(), 1)

	require.Nil(t, user)
	require.Error(t, err)
	assert.EqualError(t, err, "db timeout")
	mockRepo.AssertExpectations(t)
}

// ─────────────────────────────────────────────
// Delete
// ─────────────────────────────────────────────

func TestUserService_Delete_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewAuthService(mockRepo)

	mockRepo.On("Delete", context.Background(), int64(1)).Return(nil)

	err := svc.Delete(context.Background(), 1)

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Delete_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewAuthService(mockRepo)

	notFound := util.NewNotFoundError("no user found with id 99")
	mockRepo.On("Delete", context.Background(), int64(99)).Return(notFound)

	err := svc.Delete(context.Background(), 99)

	require.Error(t, err)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok)
	assert.Equal(t, 404, appErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Delete_InternalError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewAuthService(mockRepo)

	internalErr := util.NewInternalError("failed to delete user")
	mockRepo.On("Delete", context.Background(), int64(1)).Return(internalErr)

	err := svc.Delete(context.Background(), 1)

	require.Error(t, err)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok)
	assert.Equal(t, 500, appErr.Code)
	mockRepo.AssertExpectations(t)
}

// ─────────────────────────────────────────────
// GoogleLogin
// ─────────────────────────────────────────────

func TestUserService_GoogleLogin_Success_ExistingUser(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewAuthService(mockRepo)

	googleID := "google-sub-12345"
	existingUser := &model.User{
		ID:       42,
		Name:     "Alice",
		Email:    "alice@google.com",
		GoogleID: &googleID,
	}
	userInfo := &dto.GoogleUserInfo{
		Sub:           "google-sub-12345",
		Name:          "Alice",
		Email:         "alice@google.com",
		EmailVerified: true,
	}

	mockRepo.On("GetByGoogleID", context.Background(), "google-sub-12345").Return(existingUser, nil)

	resp, err := svc.GoogleLogin(context.Background(), userInfo)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, int64(42), resp.ID)
	assert.Equal(t, "alice@google.com", resp.Email)
	assert.NotEmpty(t, resp.Token)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "CreateGoogleUser")
}

func TestUserService_GoogleLogin_Success_NewUser(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewAuthService(mockRepo)

	userInfo := &dto.GoogleUserInfo{
		Sub:           "google-sub-99999",
		Name:          "Bob",
		Email:         "bob@google.com",
		EmailVerified: true,
	}

	mockRepo.On("GetByGoogleID", context.Background(), "google-sub-99999").Return(nil, nil)
	mockRepo.On("CreateGoogleUser", context.Background(), mock.AnythingOfType("*model.User")).Return(nil)

	resp, err := svc.GoogleLogin(context.Background(), userInfo)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "bob@google.com", resp.Email)
	assert.NotEmpty(t, resp.Token)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GoogleLogin_GetByGoogleIDError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewAuthService(mockRepo)

	userInfo := &dto.GoogleUserInfo{
		Sub:   "google-sub-12345",
		Name:  "Alice",
		Email: "alice@google.com",
	}

	mockRepo.On("GetByGoogleID", context.Background(), "google-sub-12345").
		Return(nil, errors.New("db connection refused"))

	resp, err := svc.GoogleLogin(context.Background(), userInfo)

	require.Nil(t, resp)
	require.Error(t, err)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok)
	assert.Equal(t, 500, appErr.Code)
	assert.Equal(t, "failed to fetch user", appErr.Message)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GoogleLogin_CreateGoogleUserError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewAuthService(mockRepo)

	userInfo := &dto.GoogleUserInfo{
		Sub:   "google-sub-99999",
		Name:  "Bob",
		Email: "bob@google.com",
	}

	mockRepo.On("GetByGoogleID", context.Background(), "google-sub-99999").Return(nil, nil)
	mockRepo.On("CreateGoogleUser", context.Background(), mock.AnythingOfType("*model.User")).
		Return(errors.New("db failure"))

	resp, err := svc.GoogleLogin(context.Background(), userInfo)

	require.Nil(t, resp)
	require.Error(t, err)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok)
	assert.Equal(t, 500, appErr.Code)
	assert.Equal(t, "failed to create user", appErr.Message)
	mockRepo.AssertExpectations(t)
}
