package service

import (
	"context"

	dto "gin-investment-tracker/internal/dtos"
	repository "gin-investment-tracker/internal/repositories"
)

type DashboardService struct {
	repo repository.DashboardRepositoryInterface
}

func NewDashboardService(repo repository.DashboardRepositoryInterface) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetDashboardData(ctx context.Context, userID int64) (*dto.DashboardDataDto, error) {
	data, err := s.repo.GetDashboardData(ctx, userID)
	if err != nil {
		return nil, err
	}

	data.CurrentInvestmentValue = round2(data.CurrentInvestmentValue)
	data.PreviousDayInvestmentValue = round2(data.PreviousDayInvestmentValue)
	data.TotalInvestedValue = round2(data.TotalInvestedValue)

	calculateProfitDetails(&data.TopHoldings)

	for i := range data.RecentTransactions {
		data.RecentTransactions[i].Price = round2(data.RecentTransactions[i].Price)
	}

	return data, nil
}
