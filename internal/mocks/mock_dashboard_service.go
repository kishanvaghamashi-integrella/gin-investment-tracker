package mocks

import (
	"context"

	dto "gin-investment-tracker/internal/dtos"

	"github.com/stretchr/testify/mock"
)

type MockDashboardService struct {
	mock.Mock
}

func (m *MockDashboardService) GetDashboardData(ctx context.Context, userID int64) (*dto.DashboardDataDto, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DashboardDataDto), args.Error(1)
}
