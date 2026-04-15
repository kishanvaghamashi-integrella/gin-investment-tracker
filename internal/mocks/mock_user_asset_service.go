package mocks

import (
	"context"

	dto "gin-investment-tracker/internal/dtos"
	model "gin-investment-tracker/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockUserAssetService struct {
	mock.Mock
}

func (m *MockUserAssetService) Create(ctx context.Context, userID int64, req *dto.CreateUserAssetRequest) (*model.UserAsset, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserAsset), args.Error(1)
}

func (m *MockUserAssetService) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]model.UserAsset, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.UserAsset), args.Error(1)
}

func (m *MockUserAssetService) Delete(ctx context.Context, userID, userAssetID int64) error {
	args := m.Called(ctx, userID, userAssetID)
	return args.Error(0)
}
