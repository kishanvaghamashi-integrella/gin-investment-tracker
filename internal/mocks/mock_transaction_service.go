package mocks

import (
	"context"

	dto "gin-investment-tracker/internal/dtos"
	model "gin-investment-tracker/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) Create(ctx context.Context, req *dto.CreateTransactionRequest, userId int64) (*model.Transaction, error) {
	args := m.Called(ctx, req, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Transaction), args.Error(1)
}

func (m *MockTransactionService) GetAllByUserID(ctx context.Context, userID int64, limit, offset int) ([]dto.ResponseTransactionDto, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.ResponseTransactionDto), args.Error(1)
}

func (m *MockTransactionService) Update(ctx context.Context, id int64, req *dto.UpdateTransactionRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *MockTransactionService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
