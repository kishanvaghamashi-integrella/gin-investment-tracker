package service_test

import (
	"context"
	dto "gin-investment-tracker/internal/dtos"
	"gin-investment-tracker/internal/mocks"
	model "gin-investment-tracker/internal/models"
	service "gin-investment-tracker/internal/services"
	"gin-investment-tracker/internal/util"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ─────────────────────────────────────────────
// Create
// ─────────────────────────────────────────────

func TestAssetService_Create_Success_DefaultCurrency(t *testing.T) {
	mockRepo := new(mocks.MockAssetRepository)
	svc := service.NewAssetService(mockRepo)

	req := &dto.CreateAssetRequest{
		Symbol:         "INFY",
		Name:           "Infosys Ltd",
		InstrumentType: "stock",
		ISIN:           "INE009A01021",
		Exchange:       "NSE",
	}

	mockRepo.On("Create", context.Background(), mock.AnythingOfType("*model.Asset")).Return(nil)

	asset, err := svc.Create(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, asset)
	assert.Equal(t, "INFY", asset.Symbol)
	assert.Equal(t, "INR", asset.Currency, "currency should default to INR when not provided")
	mockRepo.AssertExpectations(t)
}

func TestAssetService_Create_Success_WithExplicitCurrency(t *testing.T) {
	mockRepo := new(mocks.MockAssetRepository)
	svc := service.NewAssetService(mockRepo)

	req := &dto.CreateAssetRequest{
		Symbol:         "AAPL",
		Name:           "Apple Inc.",
		InstrumentType: "stock",
		ISIN:           "US0378331005",
		Exchange:       "NASDAQ",
		Currency:       "USD",
	}

	mockRepo.On("Create", context.Background(), mock.AnythingOfType("*model.Asset")).Return(nil)

	asset, err := svc.Create(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, asset)
	assert.Equal(t, "USD", asset.Currency, "explicit currency should be preserved")
	mockRepo.AssertExpectations(t)
}

func TestAssetService_Create_RepoError(t *testing.T) {
	mockRepo := new(mocks.MockAssetRepository)
	svc := service.NewAssetService(mockRepo)

	req := &dto.CreateAssetRequest{
		Symbol:         "INFY",
		Name:           "Infosys Ltd",
		InstrumentType: "stock",
		ISIN:           "INE009A01021",
		Exchange:       "NSE",
	}

	mockRepo.On("Create", context.Background(), mock.AnythingOfType("*model.Asset")).
		Return(util.NewInternalError("db failure"))

	asset, err := svc.Create(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, asset)
	assert.Equal(t, "db failure", err.Error())
	mockRepo.AssertExpectations(t)
}

// ─────────────────────────────────────────────
// GetByID
// ─────────────────────────────────────────────

func TestAssetService_GetByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockAssetRepository)
	svc := service.NewAssetService(mockRepo)

	expected := &model.Asset{ID: 1, Symbol: "INFY", Name: "Infosys Ltd"}
	mockRepo.On("GetByID", context.Background(), int64(1)).Return(expected, nil)

	asset, err := svc.GetByID(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, expected, asset)
	mockRepo.AssertExpectations(t)
}

func TestAssetService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockAssetRepository)
	svc := service.NewAssetService(mockRepo)

	mockRepo.On("GetByID", context.Background(), int64(99)).
		Return(nil, util.NewNotFoundError("asset not found"))

	asset, err := svc.GetByID(context.Background(), 99)

	require.Error(t, err)
	assert.Nil(t, asset)
	assert.Equal(t, "asset not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestAssetService_GetByID_InternalError(t *testing.T) {
	mockRepo := new(mocks.MockAssetRepository)
	svc := service.NewAssetService(mockRepo)

	mockRepo.On("GetByID", context.Background(), int64(1)).
		Return(nil, util.NewInternalError("db failure"))

	asset, err := svc.GetByID(context.Background(), 1)

	require.Error(t, err)
	assert.Nil(t, asset)
	mockRepo.AssertExpectations(t)
}

// ─────────────────────────────────────────────
// GetAll
// ─────────────────────────────────────────────

func TestAssetService_GetAll_Success(t *testing.T) {
	mockRepo := new(mocks.MockAssetRepository)
	svc := service.NewAssetService(mockRepo)

	expected := []model.Asset{
		{ID: 1, Symbol: "INFY"},
		{ID: 2, Symbol: "TCS"},
	}
	mockRepo.On("GetAll", context.Background(), 50, 0).Return(expected, nil)

	assets, err := svc.GetAll(context.Background(), 50, 0)

	require.NoError(t, err)
	assert.Len(t, assets, 2)
	assert.Equal(t, expected, assets)
	mockRepo.AssertExpectations(t)
}

func TestAssetService_GetAll_InternalError(t *testing.T) {
	mockRepo := new(mocks.MockAssetRepository)
	svc := service.NewAssetService(mockRepo)

	mockRepo.On("GetAll", context.Background(), 50, 0).
		Return(nil, util.NewInternalError("db failure"))

	assets, err := svc.GetAll(context.Background(), 50, 0)

	require.Error(t, err)
	assert.Nil(t, assets)
	mockRepo.AssertExpectations(t)
}

// ─────────────────────────────────────────────
// Update
// ─────────────────────────────────────────────

func TestAssetService_Update_Success(t *testing.T) {
	mockRepo := new(mocks.MockAssetRepository)
	svc := service.NewAssetService(mockRepo)

	existing := &model.Asset{ID: 1, Symbol: "INFY", Name: "Infosys Ltd", Currency: "INR"}
	newName := "Infosys Limited"
	req := &dto.UpdateAssetRequest{Name: &newName}

	mockRepo.On("GetByID", context.Background(), int64(1)).Return(existing, nil)
	mockRepo.On("Update", context.Background(), mock.AnythingOfType("*model.Asset")).Return(nil)

	err := svc.Update(context.Background(), 1, req)

	require.NoError(t, err)
	assert.Equal(t, "Infosys Limited", existing.Name, "name should have been patched in-place")
	assert.Equal(t, "INFY", existing.Symbol, "untouched fields should remain unchanged")
	mockRepo.AssertExpectations(t)
}

func TestAssetService_Update_GetByIDNotFound(t *testing.T) {
	mockRepo := new(mocks.MockAssetRepository)
	svc := service.NewAssetService(mockRepo)

	newName := "Infosys Limited"
	req := &dto.UpdateAssetRequest{Name: &newName}

	mockRepo.On("GetByID", context.Background(), int64(99)).
		Return(nil, util.NewNotFoundError("asset not found"))

	err := svc.Update(context.Background(), 99, req)

	require.Error(t, err)
	assert.Equal(t, "asset not found", err.Error())
	mockRepo.AssertNotCalled(t, "Update")
	mockRepo.AssertExpectations(t)
}

func TestAssetService_Update_RepoUpdateError(t *testing.T) {
	mockRepo := new(mocks.MockAssetRepository)
	svc := service.NewAssetService(mockRepo)

	existing := &model.Asset{ID: 1, Symbol: "INFY", Name: "Infosys Ltd"}
	newName := "Infosys Limited"
	req := &dto.UpdateAssetRequest{Name: &newName}

	mockRepo.On("GetByID", context.Background(), int64(1)).Return(existing, nil)
	mockRepo.On("Update", context.Background(), mock.AnythingOfType("*model.Asset")).
		Return(util.NewInternalError("db failure"))

	err := svc.Update(context.Background(), 1, req)

	require.Error(t, err)
	assert.Equal(t, "db failure", err.Error())
	mockRepo.AssertExpectations(t)
}

// ─────────────────────────────────────────────
// Delete
// ─────────────────────────────────────────────

func TestAssetService_Delete_Success(t *testing.T) {
	mockRepo := new(mocks.MockAssetRepository)
	svc := service.NewAssetService(mockRepo)

	mockRepo.On("Delete", context.Background(), int64(1)).Return(nil)

	err := svc.Delete(context.Background(), 1)

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAssetService_Delete_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockAssetRepository)
	svc := service.NewAssetService(mockRepo)

	mockRepo.On("Delete", context.Background(), int64(99)).
		Return(util.NewNotFoundError("asset not found"))

	err := svc.Delete(context.Background(), 99)

	require.Error(t, err)
	assert.Equal(t, "asset not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestAssetService_Delete_InternalError(t *testing.T) {
	mockRepo := new(mocks.MockAssetRepository)
	svc := service.NewAssetService(mockRepo)

	mockRepo.On("Delete", context.Background(), int64(1)).
		Return(util.NewInternalError("db failure"))

	err := svc.Delete(context.Background(), 1)

	require.Error(t, err)
	mockRepo.AssertExpectations(t)
}
