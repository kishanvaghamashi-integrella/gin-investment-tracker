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

	if data.CurrentInvestmentValue != nil {
		v1 := round2(*data.CurrentInvestmentValue)
		data.CurrentInvestmentValue = &v1
	}
	if data.PreviousDayInvestmentValue != nil {
		v2 := round2(*data.PreviousDayInvestmentValue)
		data.PreviousDayInvestmentValue = &v2
	}
	if data.TotalInvestedValue != nil {
		v3 := round2(*data.TotalInvestedValue)
		data.TotalInvestedValue = &v3
	}

	calculateProfitDetails(&data.TopHoldings)

	for i := range data.RecentTransactions {
		data.RecentTransactions[i].Price = round2(data.RecentTransactions[i].Price)
	}

	return data, nil
}
