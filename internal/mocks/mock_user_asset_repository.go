package mocks

import (
	"context"

	model "gin-investment-tracker/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockUserAssetRepository struct {
	mock.Mock
}

func (m *MockUserAssetRepository) Create(ctx context.Context, userAsset *model.UserAsset) error {
	args := m.Called(ctx, userAsset)
	return args.Error(0)
}

func (m *MockUserAssetRepository) GetIdByUserIdAssetId(ctx context.Context, userID, assetID int64) (*int64, error) {
	args := m.Called(ctx, userID, assetID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int64), args.Error(1)
}

func (m *MockUserAssetRepository) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]model.UserAsset, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.UserAsset), args.Error(1)
}

func (m *MockUserAssetRepository) Delete(ctx context.Context, id, userID int64) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *MockUserAssetRepository) IsUserAssetExits(ctx context.Context, userID int64, assetID int64) (bool, error) {
	args := m.Called(ctx, userID, assetID)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserAssetRepository) ExistsByID(ctx context.Context, id int64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}
