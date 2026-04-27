package mocks

import (
	"context"
	model "gin-investment-tracker/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockAssetRepository struct {
	mock.Mock
}

func (m *MockAssetRepository) Create(ctx context.Context, asset *model.Asset) error {
	args := m.Called(ctx, asset)
	return args.Error(0)
}

func (m *MockAssetRepository) GetByID(ctx context.Context, id int64) (*model.Asset, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Asset), args.Error(1)
}

func (m *MockAssetRepository) GetByISIN(ctx context.Context, isin string) (*model.Asset, error) {
	args := m.Called(ctx, isin)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Asset), args.Error(1)
}

func (m *MockAssetRepository) GetAll(ctx context.Context, limit, offset int) ([]model.Asset, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Asset), args.Error(1)
}

func (m *MockAssetRepository) Update(ctx context.Context, asset *model.Asset) error {
	args := m.Called(ctx, asset)
	return args.Error(0)
}

func (m *MockAssetRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAssetRepository) ExistsByID(ctx context.Context, id int64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}
