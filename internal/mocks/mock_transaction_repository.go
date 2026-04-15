package mocks

import (
	"context"

	dto "gin-investment-tracker/internal/dtos"
	model "gin-investment-tracker/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(ctx context.Context, txn *model.Transaction, holding *model.Holding, isUpdate bool) error {
	args := m.Called(ctx, txn, holding, isUpdate)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetAllByUserID(ctx context.Context, userID int64, limit, offset int) ([]dto.ResponseTransactionDto, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.ResponseTransactionDto), args.Error(1)
}

func (m *MockTransactionRepository) GetHoldingsByUserAssetID(ctx context.Context, userAssetID int64) (*model.Holding, error) {
	args := m.Called(ctx, userAssetID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Holding), args.Error(1)
}

func (m *MockTransactionRepository) GetByID(ctx context.Context, id int64) (*model.Transaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Update(ctx context.Context, txn *model.Transaction) error {
	args := m.Called(ctx, txn)
	return args.Error(0)
}

func (m *MockTransactionRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
