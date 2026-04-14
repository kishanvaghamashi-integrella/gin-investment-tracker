package mocks

import (
	"context"
	dto "gin-investment-tracker/internal/dtos"
	model "gin-investment-tracker/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(ctx context.Context, req *dto.CreateUserRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockUserService) Login(ctx context.Context, req *dto.LoginRequest) (*model.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) Delete(ctx context.Context, userId int64) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}
