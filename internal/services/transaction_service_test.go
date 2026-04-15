package service_test

import (
	"context"
	"time"

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

// helpers

func ptr[T any](v T) *T { return &v }

func newTransactionService(
	txnRepo *mocks.MockTransactionRepository,
	userAssetRepo *mocks.MockUserAssetRepository,
	userRepo *mocks.MockUserRepository,
	assetRepo *mocks.MockAssetRepository,
) *service.TransactionService {
	return service.NewTransactionService(txnRepo, userAssetRepo, userRepo, assetRepo)
}

func validCreateReq() *dto.CreateTransactionRequest {
	return &dto.CreateTransactionRequest{
		AssetID:  1,
		TxnType:  "BUY",
		Quantity: 10,
		Price:    100,
		TxnDate:  time.Now(),
	}
}

// ─────────────────────────────────────────────
// Create
// ─────────────────────────────────────────────

func TestTransactionService_Create_Success_NewHolding(t *testing.T) {
	txnRepo := new(mocks.MockTransactionRepository)
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := newTransactionService(txnRepo, userAssetRepo, userRepo, assetRepo)

	req := validCreateReq()
	userAssetID := int64(10)

	userRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	assetRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	userAssetRepo.On("GetIdByUserIdAssetId", context.Background(), int64(1), int64(1)).Return(&userAssetID, nil)
	txnRepo.On("GetHoldingsByUserAssetID", context.Background(), userAssetID).Return(nil, nil)
	txnRepo.On("Create", context.Background(), mock.AnythingOfType("*model.Transaction"), mock.AnythingOfType("*model.Holding"), false).Return(nil)

	txn, err := svc.Create(context.Background(), req, 1)

	require.NoError(t, err)
	require.NotNil(t, txn)
	assert.Equal(t, "BUY", txn.TxnType)
	assert.Equal(t, float64(10), txn.Quantity)
	userRepo.AssertExpectations(t)
	assetRepo.AssertExpectations(t)
	userAssetRepo.AssertExpectations(t)
	txnRepo.AssertExpectations(t)
}

func TestTransactionService_Create_Success_ExistingHolding_BUY(t *testing.T) {
	txnRepo := new(mocks.MockTransactionRepository)
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := newTransactionService(txnRepo, userAssetRepo, userRepo, assetRepo)

	req := validCreateReq() // BUY 10 @ 100
	userAssetID := int64(10)
	existingHolding := &model.Holding{
		UserAssetID:   userAssetID,
		TotalQuantity: 5,
		AveragePrice:  80,
		TotalInvested: 400,
	}

	userRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	assetRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	userAssetRepo.On("GetIdByUserIdAssetId", context.Background(), int64(1), int64(1)).Return(&userAssetID, nil)
	txnRepo.On("GetHoldingsByUserAssetID", context.Background(), userAssetID).Return(existingHolding, nil)
	txnRepo.On("Create", context.Background(), mock.AnythingOfType("*model.Transaction"), mock.AnythingOfType("*model.Holding"), true).Return(nil)

	txn, err := svc.Create(context.Background(), req, 1)

	require.NoError(t, err)
	require.NotNil(t, txn)
	// Holding math: (5*80 + 10*100) / 15 = (400+1000)/15 ≈ 93.33
	assert.Equal(t, float64(15), existingHolding.TotalQuantity)
	assert.InDelta(t, 93.33, existingHolding.AveragePrice, 0.01)
	assert.Equal(t, float64(1400), existingHolding.TotalInvested)
	txnRepo.AssertExpectations(t)
}

func TestTransactionService_Create_Success_ExistingHolding_SELL(t *testing.T) {
	txnRepo := new(mocks.MockTransactionRepository)
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := newTransactionService(txnRepo, userAssetRepo, userRepo, assetRepo)

	req := &dto.CreateTransactionRequest{
		AssetID:  1,
		TxnType:  "SELL",
		Quantity: 3,
		Price:    110,
		TxnDate:  time.Now(),
	}
	userAssetID := int64(10)
	existingHolding := &model.Holding{
		UserAssetID:   userAssetID,
		TotalQuantity: 5,
		AveragePrice:  80,
		TotalInvested: 400,
	}

	userRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	assetRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	userAssetRepo.On("GetIdByUserIdAssetId", context.Background(), int64(1), int64(1)).Return(&userAssetID, nil)
	txnRepo.On("GetHoldingsByUserAssetID", context.Background(), userAssetID).Return(existingHolding, nil)
	txnRepo.On("Create", context.Background(), mock.AnythingOfType("*model.Transaction"), mock.AnythingOfType("*model.Holding"), true).Return(nil)

	txn, err := svc.Create(context.Background(), req, 1)

	require.NoError(t, err)
	require.NotNil(t, txn)
	// Holding after sell: qty=2, invested=400-(110*3)=70
	assert.Equal(t, float64(2), existingHolding.TotalQuantity)
	assert.Equal(t, float64(70), existingHolding.TotalInvested)
	txnRepo.AssertExpectations(t)
}

func TestTransactionService_Create_Success_NoUserAsset_CreatesOne(t *testing.T) {
	txnRepo := new(mocks.MockTransactionRepository)
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := newTransactionService(txnRepo, userAssetRepo, userRepo, assetRepo)

	req := validCreateReq()

	userRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	assetRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	userAssetRepo.On("GetIdByUserIdAssetId", context.Background(), int64(1), int64(1)).Return(nil, nil)
	userAssetRepo.On("Create", context.Background(), mock.AnythingOfType("*model.UserAsset")).Return(nil)
	txnRepo.On("GetHoldingsByUserAssetID", context.Background(), int64(0)).Return(nil, nil)
	txnRepo.On("Create", context.Background(), mock.AnythingOfType("*model.Transaction"), mock.AnythingOfType("*model.Holding"), false).Return(nil)

	txn, err := svc.Create(context.Background(), req, 1)

	require.NoError(t, err)
	require.NotNil(t, txn)
	userAssetRepo.AssertCalled(t, "Create", context.Background(), mock.AnythingOfType("*model.UserAsset"))
}

func TestTransactionService_Create_UserNotFound(t *testing.T) {
	txnRepo := new(mocks.MockTransactionRepository)
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := newTransactionService(txnRepo, userAssetRepo, userRepo, assetRepo)

	userRepo.On("ExistsByID", context.Background(), int64(99)).Return(false, nil)

	txn, err := svc.Create(context.Background(), validCreateReq(), 99)

	require.Nil(t, txn)
	require.Error(t, err)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok)
	assert.Equal(t, 404, appErr.Code)
	assert.Contains(t, appErr.Message, "user not found")
	userRepo.AssertExpectations(t)
	txnRepo.AssertNotCalled(t, "Create")
}

func TestTransactionService_Create_UserRepoError(t *testing.T) {
	txnRepo := new(mocks.MockTransactionRepository)
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := newTransactionService(txnRepo, userAssetRepo, userRepo, assetRepo)

	userRepo.On("ExistsByID", context.Background(), int64(1)).Return(false, util.NewInternalError("db failure"))

	txn, err := svc.Create(context.Background(), validCreateReq(), 1)

	require.Nil(t, txn)
	require.Error(t, err)
	txnRepo.AssertNotCalled(t, "Create")
}

func TestTransactionService_Create_AssetNotFound(t *testing.T) {
	txnRepo := new(mocks.MockTransactionRepository)
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := newTransactionService(txnRepo, userAssetRepo, userRepo, assetRepo)

	userRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	assetRepo.On("ExistsByID", context.Background(), int64(1)).Return(false, nil)

	txn, err := svc.Create(context.Background(), validCreateReq(), 1)

	require.Nil(t, txn)
	require.Error(t, err)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok)
	assert.Equal(t, 404, appErr.Code)
	assert.Contains(t, appErr.Message, "asset not found")
	txnRepo.AssertNotCalled(t, "Create")
}

func TestTransactionService_Create_SellWithNoHolding(t *testing.T) {
	txnRepo := new(mocks.MockTransactionRepository)
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := newTransactionService(txnRepo, userAssetRepo, userRepo, assetRepo)

	req := &dto.CreateTransactionRequest{
		AssetID:  1,
		TxnType:  "SELL",
		Quantity: 5,
		Price:    100,
		TxnDate:  time.Now(),
	}
	userAssetID := int64(10)

	userRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	assetRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	userAssetRepo.On("GetIdByUserIdAssetId", context.Background(), int64(1), int64(1)).Return(&userAssetID, nil)
	txnRepo.On("GetHoldingsByUserAssetID", context.Background(), userAssetID).Return(nil, nil)

	txn, err := svc.Create(context.Background(), req, 1)

	require.Nil(t, txn)
	require.Error(t, err)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok)
	assert.Equal(t, 400, appErr.Code)
	assert.Contains(t, appErr.Message, "Cannot sell asset that is not currently held")
	txnRepo.AssertNotCalled(t, "Create")
}

func TestTransactionService_Create_SellExceedsHolding(t *testing.T) {
	txnRepo := new(mocks.MockTransactionRepository)
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := newTransactionService(txnRepo, userAssetRepo, userRepo, assetRepo)

	req := &dto.CreateTransactionRequest{
		AssetID:  1,
		TxnType:  "SELL",
		Quantity: 20,
		Price:    100,
		TxnDate:  time.Now(),
	}
	userAssetID := int64(10)
	existingHolding := &model.Holding{
		UserAssetID:   userAssetID,
		TotalQuantity: 5,
		AveragePrice:  80,
		TotalInvested: 400,
	}

	userRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	assetRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	userAssetRepo.On("GetIdByUserIdAssetId", context.Background(), int64(1), int64(1)).Return(&userAssetID, nil)
	txnRepo.On("GetHoldingsByUserAssetID", context.Background(), userAssetID).Return(existingHolding, nil)

	txn, err := svc.Create(context.Background(), req, 1)

	require.Nil(t, txn)
	require.Error(t, err)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok)
	assert.Equal(t, 400, appErr.Code)
	assert.Contains(t, appErr.Message, "Sell quantity exceeds current holding quantity")
	txnRepo.AssertNotCalled(t, "Create")
}

func TestTransactionService_Create_RepoCreateError(t *testing.T) {
	txnRepo := new(mocks.MockTransactionRepository)
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := newTransactionService(txnRepo, userAssetRepo, userRepo, assetRepo)

	userAssetID := int64(10)

	userRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	assetRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	userAssetRepo.On("GetIdByUserIdAssetId", context.Background(), int64(1), int64(1)).Return(&userAssetID, nil)
	txnRepo.On("GetHoldingsByUserAssetID", context.Background(), userAssetID).Return(nil, nil)
	txnRepo.On("Create", context.Background(), mock.AnythingOfType("*model.Transaction"), mock.AnythingOfType("*model.Holding"), false).
		Return(util.NewInternalError("db failure"))

	txn, err := svc.Create(context.Background(), validCreateReq(), 1)

	require.Nil(t, txn)
	require.Error(t, err)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok)
	assert.Equal(t, 500, appErr.Code)
}

// ─────────────────────────────────────────────
// GetAllByUserID
// ─────────────────────────────────────────────

func TestTransactionService_GetAllByUserID_Success(t *testing.T) {
	txnRepo := new(mocks.MockTransactionRepository)
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := newTransactionService(txnRepo, userAssetRepo, userRepo, assetRepo)

	expected := []dto.ResponseTransactionDto{
		{ID: 1, AssetName: "INFY", TxnType: "BUY", Quantity: 10, Price: 100},
		{ID: 2, AssetName: "TCS", TxnType: "SELL", Quantity: 5, Price: 200},
	}

	userRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	txnRepo.On("GetAllByUserID", context.Background(), int64(1), 50, 0).Return(expected, nil)

	result, err := svc.GetAllByUserID(context.Background(), 1, 50, 0)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, expected, result)
	userRepo.AssertExpectations(t)
	txnRepo.AssertExpectations(t)
}

func TestTransactionService_GetAllByUserID_UserNotFound(t *testing.T) {
	txnRepo := new(mocks.MockTransactionRepository)
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := newTransactionService(txnRepo, userAssetRepo, userRepo, assetRepo)

	userRepo.On("ExistsByID", context.Background(), int64(99)).Return(false, nil)

	result, err := svc.GetAllByUserID(context.Background(), 99, 50, 0)

	require.Nil(t, result)
	require.Error(t, err)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok)
	assert.Equal(t, 404, appErr.Code)
	txnRepo.AssertNotCalled(t, "GetAllByUserID")
}

func TestTransactionService_GetAllByUserID_RepoError(t *testing.T) {
	txnRepo := new(mocks.MockTransactionRepository)
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := newTransactionService(txnRepo, userAssetRepo, userRepo, assetRepo)

	userRepo.On("ExistsByID", context.Background(), int64(1)).Return(true, nil)
	txnRepo.On("GetAllByUserID", context.Background(), int64(1), 50, 0).
		Return(nil, util.NewInternalError("db failure"))

	result, err := svc.GetAllByUserID(context.Background(), 1, 50, 0)

	require.Nil(t, result)
	require.Error(t, err)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok)
	assert.Equal(t, 500, appErr.Code)
}

// ─────────────────────────────────────────────
// Update
// ─────────────────────────────────────────────

func TestTransactionService_Update_Success(t *testing.T) {
	txnRepo := new(mocks.MockTransactionRepository)
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := newTransactionService(txnRepo, userAssetRepo, userRepo, assetRepo)

	existing := &model.Transaction{ID: 1, TxnType: "BUY", Quantity: 10, Price: 100}
	newQty := float64(20)
	req := &dto.UpdateTransactionRequest{Quantity: &newQty}

	txnRepo.On("GetByID", context.Background(), int64(1)).Return(existing, nil)
	txnRepo.On("Update", context.Background(), mock.AnythingOfType("*model.Transaction")).Return(nil)

	err := svc.Update(context.Background(), 1, req)

	require.NoError(t, err)
	assert.Equal(t, float64(20), existing.Quantity, "quantity should be patched in-place")
	assert.Equal(t, "BUY", existing.TxnType, "unpatched fields must remain unchanged")
	txnRepo.AssertExpectations(t)
}

func TestTransactionService_Update_PatchAllFields(t *testing.T) {
	txnRepo := new(mocks.MockTransactionRepository)
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := newTransactionService(txnRepo, userAssetRepo, userRepo, assetRepo)

	newDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	existing := &model.Transaction{ID: 1, TxnType: "BUY", Quantity: 10, Price: 100}
	req := &dto.UpdateTransactionRequest{
		TxnType:  ptr("SELL"),
		Quantity: ptr(5.0),
		Price:    ptr(150.0),
		TxnDate:  &newDate,
	}

	txnRepo.On("GetByID", context.Background(), int64(1)).Return(existing, nil)
	txnRepo.On("Update", context.Background(), mock.AnythingOfType("*model.Transaction")).Return(nil)

	err := svc.Update(context.Background(), 1, req)

	require.NoError(t, err)
	assert.Equal(t, "SELL", existing.TxnType)
	assert.Equal(t, float64(5), existing.Quantity)
	assert.Equal(t, float64(150), existing.Price)
	assert.Equal(t, newDate, existing.TxnDate)
	txnRepo.AssertExpectations(t)
}

func TestTransactionService_Update_NotFound(t *testing.T) {
	txnRepo := new(mocks.MockTransactionRepository)
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := newTransactionService(txnRepo, userAssetRepo, userRepo, assetRepo)

	txnRepo.On("GetByID", context.Background(), int64(99)).
		Return(nil, util.NewNotFoundError("transaction with id 99 not found"))

	err := svc.Update(context.Background(), 99, &dto.UpdateTransactionRequest{})

	require.Error(t, err)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok)
	assert.Equal(t, 404, appErr.Code)
	txnRepo.AssertNotCalled(t, "Update")
}

func TestTransactionService_Update_RepoUpdateError(t *testing.T) {
	txnRepo := new(mocks.MockTransactionRepository)
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := newTransactionService(txnRepo, userAssetRepo, userRepo, assetRepo)

	existing := &model.Transaction{ID: 1, TxnType: "BUY", Quantity: 10, Price: 100}
	newQty := float64(20)

	txnRepo.On("GetByID", context.Background(), int64(1)).Return(existing, nil)
	txnRepo.On("Update", context.Background(), mock.AnythingOfType("*model.Transaction")).
		Return(util.NewInternalError("db failure"))

	err := svc.Update(context.Background(), 1, &dto.UpdateTransactionRequest{Quantity: &newQty})

	require.Error(t, err)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok)
	assert.Equal(t, 500, appErr.Code)
	txnRepo.AssertExpectations(t)
}

// ─────────────────────────────────────────────
// Delete
// ─────────────────────────────────────────────

func TestTransactionService_Delete_Success(t *testing.T) {
	txnRepo := new(mocks.MockTransactionRepository)
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := newTransactionService(txnRepo, userAssetRepo, userRepo, assetRepo)

	txnRepo.On("Delete", context.Background(), int64(1)).Return(nil)

	err := svc.Delete(context.Background(), 1)

	require.NoError(t, err)
	txnRepo.AssertExpectations(t)
}

func TestTransactionService_Delete_NotFound(t *testing.T) {
	txnRepo := new(mocks.MockTransactionRepository)
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := newTransactionService(txnRepo, userAssetRepo, userRepo, assetRepo)

	txnRepo.On("Delete", context.Background(), int64(99)).
		Return(util.NewNotFoundError("transaction with id 99 not found"))

	err := svc.Delete(context.Background(), 99)

	require.Error(t, err)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok)
	assert.Equal(t, 404, appErr.Code)
	txnRepo.AssertExpectations(t)
}

func TestTransactionService_Delete_InternalError(t *testing.T) {
	txnRepo := new(mocks.MockTransactionRepository)
	userAssetRepo := new(mocks.MockUserAssetRepository)
	userRepo := new(mocks.MockUserRepository)
	assetRepo := new(mocks.MockAssetRepository)
	svc := newTransactionService(txnRepo, userAssetRepo, userRepo, assetRepo)

	txnRepo.On("Delete", context.Background(), int64(1)).
		Return(util.NewInternalError("db failure"))

	err := svc.Delete(context.Background(), 1)

	require.Error(t, err)
	appErr, ok := err.(*util.AppError)
	require.True(t, ok)
	assert.Equal(t, 500, appErr.Code)
	txnRepo.AssertExpectations(t)
}
