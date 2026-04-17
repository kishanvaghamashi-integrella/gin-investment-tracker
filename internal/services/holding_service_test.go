package service_test

import (
	"context"
	"fmt"
	"testing"

	dto "gin-investment-tracker/internal/dtos"
	"gin-investment-tracker/internal/mocks"
	service "gin-investment-tracker/internal/services"
	"gin-investment-tracker/internal/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newHoldingService(
	holdingRepo *mocks.MockHoldingRepository,
	userRepo *mocks.MockUserRepository,
) *service.HoldingService {
	return service.NewHoldingService(holdingRepo, userRepo)
}

// ─────────────────────────────────────────────
// GetAllByUserID
// ─────────────────────────────────────────────

func TestHoldingService_GetAllByUserID_UserNotFound(t *testing.T) {
	holdingRepo := new(mocks.MockHoldingRepository)
	userRepo := new(mocks.MockUserRepository)
	svc := newHoldingService(holdingRepo, userRepo)

	userRepo.On("ExistsByID", context.Background(), int64(99)).Return(false, nil)

	result, err := svc.GetAllByUserID(context.Background(), 99, 50, 0)

	require.Error(t, err)
	assert.Nil(t, result)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok)
	assert.Equal(t, 404, appErr.Code)
	assert.Equal(t, fmt.Sprintf("user with id %d not found", 99), appErr.Message)
	userRepo.AssertExpectations(t)
	holdingRepo.AssertNotCalled(t, "GetAllByUserID")
}

func TestHoldingService_GetAllByUserID_UserRepoError(t *testing.T) {
	holdingRepo := new(mocks.MockHoldingRepository)
	userRepo := new(mocks.MockUserRepository)
	svc := newHoldingService(holdingRepo, userRepo)

	userRepo.On("ExistsByID", context.Background(), int64(1)).Return(false, util.NewInternalError("db error"))

	result, err := svc.GetAllByUserID(context.Background(), 1, 50, 0)

	require.Error(t, err)
	assert.Nil(t, result)
	userRepo.AssertExpectations(t)
	holdingRepo.AssertNotCalled(t, "GetAllByUserID")
}

func TestHoldingService_GetAllByUserID_RepoError(t *testing.T) {
	holdingRepo := new(mocks.MockHoldingRepository)
	userRepo := new(mocks.MockUserRepository)
	svc := newHoldingService(holdingRepo, userRepo)

	userRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	holdingRepo.On("GetAllByUserID", context.Background(), int64(1), 50, 0).Return(nil, util.NewInternalError("failed to list holdings"))

	result, err := svc.GetAllByUserID(context.Background(), 1, 50, 0)

	require.Error(t, err)
	assert.Nil(t, result)
	userRepo.AssertExpectations(t)
	holdingRepo.AssertExpectations(t)
}

func TestHoldingService_GetAllByUserID_Success_EmptyResult(t *testing.T) {
	holdingRepo := new(mocks.MockHoldingRepository)
	userRepo := new(mocks.MockUserRepository)
	svc := newHoldingService(holdingRepo, userRepo)

	userRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	holdingRepo.On("GetAllByUserID", context.Background(), int64(1), 50, 0).Return([]dto.HoldingResponseDto{}, nil)

	result, err := svc.GetAllByUserID(context.Background(), 1, 50, 0)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
	userRepo.AssertExpectations(t)
	holdingRepo.AssertExpectations(t)
}

func TestHoldingService_GetAllByUserID_Success_WithPagination(t *testing.T) {
	holdingRepo := new(mocks.MockHoldingRepository)
	userRepo := new(mocks.MockUserRepository)
	svc := newHoldingService(holdingRepo, userRepo)

	userRepo.On("ExistsByID", context.Background(), int64(2)).Return(true, nil)
	holdingRepo.On("GetAllByUserID", context.Background(), int64(2), 10, 5).Return([]dto.HoldingResponseDto{}, nil)

	result, err := svc.GetAllByUserID(context.Background(), 2, 10, 5)

	require.NoError(t, err)
	assert.NotNil(t, result)
	userRepo.AssertExpectations(t)
	holdingRepo.AssertExpectations(t)
}
