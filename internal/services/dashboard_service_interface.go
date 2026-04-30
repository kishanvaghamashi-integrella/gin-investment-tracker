package service

import (
	"context"

	dto "gin-investment-tracker/internal/dtos"
)

type DashboardServiceInterface interface {
	GetDashboardData(ctx context.Context, userID int64) (*dto.DashboardDataDto, error)
}
