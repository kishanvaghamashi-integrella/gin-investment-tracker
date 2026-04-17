package mocks

import (
	"context"

	dto "gin-investment-tracker/internal/dtos"

	"github.com/stretchr/testify/mock"
)

type MockHoldingService struct {
	mock.Mock
}

func (m *MockHoldingService) GetAllByUserID(ctx context.Context, userID int64, limit, offset int) ([]dto.HoldingResponseDto, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.HoldingResponseDto), args.Error(1)
}
