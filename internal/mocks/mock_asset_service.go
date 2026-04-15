package mocks

import (
	"context"
	dto "gin-investment-tracker/internal/dtos"
	model "gin-investment-tracker/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockAssetService struct {
	mock.Mock
}

func (m *MockAssetService) Create(ctx context.Context, req *dto.CreateAssetRequest) (*model.Asset, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Asset), args.Error(1)
}

func (m *MockAssetService) GetByID(ctx context.Context, id int64) (*model.Asset, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Asset), args.Error(1)
}

func (m *MockAssetService) GetAll(ctx context.Context, limit, offset int) ([]model.Asset, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Asset), args.Error(1)
}

func (m *MockAssetService) Update(ctx context.Context, id int64, req *dto.UpdateAssetRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *MockAssetService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
