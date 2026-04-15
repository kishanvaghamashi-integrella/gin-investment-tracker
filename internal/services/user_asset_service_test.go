package service_test

import (
	"context"
	"fmt"
	"testing"

	dto "gin-investment-tracker/internal/dtos"
	"gin-investment-tracker/internal/mocks"
	model "gin-investment-tracker/internal/models"
	service "gin-investment-tracker/internal/services"
	"gin-investment-tracker/internal/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ─────────────────────────────────────────────
// Create
// ─────────────────────────────────────────────

func TestUserAssetService_Create_Success(t *testing.T) {
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := service.NewUserAssetService(userAssetRepo, userRepo, assetRepo)

	req := &dto.CreateUserAssetRequest{AssetID: 5}

	userRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	assetRepo.On("ExistsByID", context.Background(), int64(5)).Return(true, nil)
	userAssetRepo.On("IsUserAssetExists", context.Background(), int64(1), int64(5)).Return(false, nil)
	userAssetRepo.On("Create", context.Background(), mock.AnythingOfType("*model.UserAsset")).Return(nil)

	result, err := svc.Create(context.Background(), 1, req)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, int64(1), result.UserID)
	assert.Equal(t, int64(5), result.AssetID)
	userRepo.AssertExpectations(t)
	assetRepo.AssertExpectations(t)
	userAssetRepo.AssertExpectations(t)
}

func TestUserAssetService_Create_UserNotFound(t *testing.T) {
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := service.NewUserAssetService(userAssetRepo, userRepo, assetRepo)

	req := &dto.CreateUserAssetRequest{AssetID: 5}

	userRepo.On("ExistsByID", context.Background(), int64(99)).Return(false, nil)

	result, err := svc.Create(context.Background(), 99, req)

	require.Error(t, err)
	assert.Nil(t, result)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok)
	assert.Equal(t, 404, appErr.Code)
	assert.Equal(t, fmt.Sprintf("user with id %d not found", 99), appErr.Message)
	userRepo.AssertExpectations(t)
	assetRepo.AssertNotCalled(t, "ExistsByID")
	userAssetRepo.AssertNotCalled(t, "Create")
}

func TestUserAssetService_Create_AssetNotFound(t *testing.T) {
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := service.NewUserAssetService(userAssetRepo, userRepo, assetRepo)

	req := &dto.CreateUserAssetRequest{AssetID: 99}

	userRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	assetRepo.On("ExistsByID", context.Background(), int64(99)).Return(false, nil)

	result, err := svc.Create(context.Background(), 1, req)

	require.Error(t, err)
	assert.Nil(t, result)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok)
	assert.Equal(t, 404, appErr.Code)
	userRepo.AssertExpectations(t)
	assetRepo.AssertExpectations(t)
	userAssetRepo.AssertNotCalled(t, "Create")
}

func TestUserAssetService_Create_AlreadyExists(t *testing.T) {
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := service.NewUserAssetService(userAssetRepo, userRepo, assetRepo)

	req := &dto.CreateUserAssetRequest{AssetID: 5}

	userRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	assetRepo.On("ExistsByID", context.Background(), int64(5)).Return(true, nil)
	userAssetRepo.On("IsUserAssetExists", context.Background(), int64(1), int64(5)).Return(true, nil)

	result, err := svc.Create(context.Background(), 1, req)

	require.Error(t, err)
	assert.Nil(t, result)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok)
	assert.Equal(t, 400, appErr.Code)
	assert.Equal(t, "This entry already exists", appErr.Message)
	userAssetRepo.AssertNotCalled(t, "Create")
}

func TestUserAssetService_Create_RepoError(t *testing.T) {
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := service.NewUserAssetService(userAssetRepo, userRepo, assetRepo)

	req := &dto.CreateUserAssetRequest{AssetID: 5}

	userRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	assetRepo.On("ExistsByID", context.Background(), int64(5)).Return(true, nil)
	userAssetRepo.On("IsUserAssetExists", context.Background(), int64(1), int64(5)).Return(false, nil)
	userAssetRepo.On("Create", context.Background(), mock.AnythingOfType("*model.UserAsset")).
		Return(util.NewInternalError("db failure"))

	result, err := svc.Create(context.Background(), 1, req)

	require.Error(t, err)
	assert.Nil(t, result)
	userAssetRepo.AssertExpectations(t)
}

// ─────────────────────────────────────────────
// GetByUserID
// ─────────────────────────────────────────────

func TestUserAssetService_GetByUserID_Success(t *testing.T) {
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := service.NewUserAssetService(userAssetRepo, userRepo, assetRepo)

	expected := []model.UserAsset{
		{ID: 1, UserID: 1, AssetID: 5},
		{ID: 2, UserID: 1, AssetID: 7},
	}

	userRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	userAssetRepo.On("GetByUserID", context.Background(), int64(1), 50, 0).Return(expected, nil)

	result, err := svc.GetByUserID(context.Background(), 1, 50, 0)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, expected, result)
	userRepo.AssertExpectations(t)
	userAssetRepo.AssertExpectations(t)
}

func TestUserAssetService_GetByUserID_UserNotFound(t *testing.T) {
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := service.NewUserAssetService(userAssetRepo, userRepo, assetRepo)

	userRepo.On("ExistsByID", context.Background(), int64(99)).Return(false, nil)

	result, err := svc.GetByUserID(context.Background(), 99, 50, 0)

	require.Error(t, err)
	assert.Nil(t, result)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok)
	assert.Equal(t, 404, appErr.Code)
	userAssetRepo.AssertNotCalled(t, "GetByUserID")
}

func TestUserAssetService_GetByUserID_RepoError(t *testing.T) {
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := service.NewUserAssetService(userAssetRepo, userRepo, assetRepo)

	userRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	userAssetRepo.On("GetByUserID", context.Background(), int64(1), 50, 0).
		Return(nil, util.NewInternalError("db failure"))

	result, err := svc.GetByUserID(context.Background(), 1, 50, 0)

	require.Error(t, err)
	assert.Nil(t, result)
	userRepo.AssertExpectations(t)
	userAssetRepo.AssertExpectations(t)
}

// ─────────────────────────────────────────────
// Delete
// ─────────────────────────────────────────────

func TestUserAssetService_Delete_Success(t *testing.T) {
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := service.NewUserAssetService(userAssetRepo, userRepo, assetRepo)

	userRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	userAssetRepo.On("Delete", context.Background(), int64(10), int64(1)).Return(nil)

	err := svc.Delete(context.Background(), 1, 10)

	require.NoError(t, err)
	userRepo.AssertExpectations(t)
	userAssetRepo.AssertExpectations(t)
}

func TestUserAssetService_Delete_UserNotFound(t *testing.T) {
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := service.NewUserAssetService(userAssetRepo, userRepo, assetRepo)

	userRepo.On("ExistsByID", context.Background(), int64(99)).Return(false, nil)

	err := svc.Delete(context.Background(), 99, 10)

	require.Error(t, err)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok)
	assert.Equal(t, 404, appErr.Code)
	userAssetRepo.AssertNotCalled(t, "Delete")
}

func TestUserAssetService_Delete_RepoError(t *testing.T) {
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := service.NewUserAssetService(userAssetRepo, userRepo, assetRepo)

	userRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	userAssetRepo.On("Delete", context.Background(), int64(10), int64(1)).
		Return(util.NewInternalError("db failure"))

	err := svc.Delete(context.Background(), 1, 10)

	require.Error(t, err)
	userRepo.AssertExpectations(t)
	userAssetRepo.AssertExpectations(t)
}
