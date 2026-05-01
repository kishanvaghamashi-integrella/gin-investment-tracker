package repository

import (
	"context"
	dto "gin-investment-tracker/internal/dtos"
)

type DashboardRepositoryInterface interface {
	GetDashboardData(ctx context.Context, userID int64) (*dto.DashboardDataDto, error)
}
