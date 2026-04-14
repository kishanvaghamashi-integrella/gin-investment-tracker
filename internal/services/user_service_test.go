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
	svc := service.NewUserService(mockRepo)

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
	svc := service.NewUserService(mockRepo)

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
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewUserService(mockRepo)

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

	user, err := svc.Login(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, storedUser.ID, user.ID)
	assert.Equal(t, storedUser.Email, user.Email)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_EmailNotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewUserService(mockRepo)

	req := &dto.LoginRequest{
		Email:    "ghost@example.com",
		Password: "secret123",
	}

	mockRepo.On("GetByEmail", context.Background(), req.Email).Return(nil, pgx.ErrNoRows)

	user, err := svc.Login(context.Background(), req)

	require.Nil(t, user)
	require.Error(t, err)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok, "expected *util.AppError")
	assert.Equal(t, 400, appErr.Code)
	assert.Equal(t, "invalid email or password", appErr.Message)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_WrongPassword(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewUserService(mockRepo)

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

	user, err := svc.Login(context.Background(), req)

	require.Nil(t, user)
	require.Error(t, err)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok, "expected *util.AppError")
	assert.Equal(t, 400, appErr.Code)
	assert.Equal(t, "invalid email or password", appErr.Message)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_RepoInternalError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewUserService(mockRepo)

	req := &dto.LoginRequest{
		Email:    "alice@example.com",
		Password: "secret123",
	}

	mockRepo.On("GetByEmail", context.Background(), req.Email).Return(nil, errors.New("connection refused"))

	user, err := svc.Login(context.Background(), req)

	require.Nil(t, user)
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
	svc := service.NewUserService(mockRepo)

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
	svc := service.NewUserService(mockRepo)

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
	svc := service.NewUserService(mockRepo)

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
	svc := service.NewUserService(mockRepo)

	mockRepo.On("Delete", context.Background(), int64(1)).Return(nil)

	err := svc.Delete(context.Background(), 1)

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Delete_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewUserService(mockRepo)

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
	svc := service.NewUserService(mockRepo)

	internalErr := util.NewInternalError("failed to delete user")
	mockRepo.On("Delete", context.Background(), int64(1)).Return(internalErr)

	err := svc.Delete(context.Background(), 1)

	require.Error(t, err)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok)
	assert.Equal(t, 500, appErr.Code)
	mockRepo.AssertExpectations(t)
}
