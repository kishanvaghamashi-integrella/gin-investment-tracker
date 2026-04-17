package service_test

import (
	"context"
	"fmt"
	"sort"
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

func newHoldingServiceWithFetcher(
	holdingRepo *mocks.MockHoldingRepository,
	userRepo *mocks.MockUserRepository,
	fetcher *mocks.MockPriceFetcher,
) *service.HoldingService {
	return service.NewHoldingServiceWithFetcher(holdingRepo, userRepo, fetcher)
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

// ─────────────────────────────────────────────
// Price calculation (deterministic via injected fetcher)
// ─────────────────────────────────────────────

func TestHoldingService_GetAllByUserID_CalculatesCurrentPrice(t *testing.T) {
	holdingRepo := new(mocks.MockHoldingRepository)
	userRepo := new(mocks.MockUserRepository)
	fetcher := new(mocks.MockPriceFetcher)
	svc := newHoldingServiceWithFetcher(holdingRepo, userRepo, fetcher)

	storedHolding := dto.HoldingResponseDto{
		ID:                      1,
		AssetExternalPlatformID: "123456",
		AssetInstrumentType:     "mutual_fund",
		Quantity:                10,
		AveragePrice:            100,
		InvestedCapital:         1000,
	}

	userRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	holdingRepo.On("GetAllByUserID", context.Background(), int64(1), 50, 0).
		Return([]dto.HoldingResponseDto{storedHolding}, nil)
	fetcher.On("FetchPrice", "mutual_fund", "123456").Return(120.0, nil)

	result, err := svc.GetAllByUserID(context.Background(), 1, 50, 0)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, 120.0, result[0].CurrentPrice)
	assert.Equal(t, 1200.0, result[0].CurrentCapital) // 120 * 10
	assert.Equal(t, 20.0, result[0].ReturnPercentage) // (1200-1000)/1000 * 100
	fetcher.AssertExpectations(t)
}

func TestHoldingService_GetAllByUserID_FallsBackToAveragePriceOnFetchError(t *testing.T) {
	holdingRepo := new(mocks.MockHoldingRepository)
	userRepo := new(mocks.MockUserRepository)
	fetcher := new(mocks.MockPriceFetcher)
	svc := newHoldingServiceWithFetcher(holdingRepo, userRepo, fetcher)

	storedHolding := dto.HoldingResponseDto{
		ID:                      2,
		AssetExternalPlatformID: "RELIANCE.NS",
		AssetInstrumentType:     "stock",
		Quantity:                5,
		AveragePrice:            200,
		InvestedCapital:         1000,
	}

	userRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	holdingRepo.On("GetAllByUserID", context.Background(), int64(1), 50, 0).
		Return([]dto.HoldingResponseDto{storedHolding}, nil)
	fetcher.On("FetchPrice", "stock", "RELIANCE.NS").Return(0.0, fmt.Errorf("api unreachable"))

	result, err := svc.GetAllByUserID(context.Background(), 1, 50, 0)

	require.NoError(t, err)
	require.Len(t, result, 1)
	// fallback: CurrentPrice should equal AveragePrice
	assert.Equal(t, 200.0, result[0].CurrentPrice)
	assert.Equal(t, 1000.0, result[0].CurrentCapital) // 200 * 5
	assert.Equal(t, 0.0, result[0].ReturnPercentage)  // (1000-1000)/1000 * 100
	fetcher.AssertExpectations(t)
}

func TestHoldingService_GetAllByUserID_MultipleHoldingsConcurrently(t *testing.T) {
	holdingRepo := new(mocks.MockHoldingRepository)
	userRepo := new(mocks.MockUserRepository)
	fetcher := new(mocks.MockPriceFetcher)
	svc := newHoldingServiceWithFetcher(holdingRepo, userRepo, fetcher)

	stored := []dto.HoldingResponseDto{
		{ID: 1, AssetExternalPlatformID: "111111", AssetInstrumentType: "mutual_fund", Quantity: 10, AveragePrice: 100, InvestedCapital: 1000},
		{ID: 2, AssetExternalPlatformID: "TCS.NS", AssetInstrumentType: "stock", Quantity: 2, AveragePrice: 3000, InvestedCapital: 6000},
	}

	userRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	holdingRepo.On("GetAllByUserID", context.Background(), int64(1), 50, 0).Return(stored, nil)
	fetcher.On("FetchPrice", "mutual_fund", "111111").Return(110.0, nil)
	fetcher.On("FetchPrice", "stock", "TCS.NS").Return(3600.0, nil)

	result, err := svc.GetAllByUserID(context.Background(), 1, 50, 0)

	require.NoError(t, err)
	require.Len(t, result, 2)

	// sort by ID for deterministic assertions (goroutines may complete in any order)
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })

	// MF holding: price=110, qty=10, invested=1000
	assert.Equal(t, 110.0, result[0].CurrentPrice)
	assert.Equal(t, 1100.0, result[0].CurrentCapital)
	assert.Equal(t, 10.0, result[0].ReturnPercentage) // (1100-1000)/1000*100

	// stock holding: price=3600, qty=2, invested=6000
	assert.Equal(t, 3600.0, result[1].CurrentPrice)
	assert.Equal(t, 7200.0, result[1].CurrentCapital)
	assert.Equal(t, 20.0, result[1].ReturnPercentage) // (7200-6000)/6000*100

	fetcher.AssertExpectations(t)
}
