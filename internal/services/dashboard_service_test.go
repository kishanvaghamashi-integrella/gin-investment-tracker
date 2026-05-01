package service_test

import (
	"context"
	"testing"
	"time"

	dto "gin-investment-tracker/internal/dtos"
	"gin-investment-tracker/internal/mocks"
	service "gin-investment-tracker/internal/services"
	"gin-investment-tracker/internal/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newDashboardService(repo *mocks.MockDashboardRepository) *service.DashboardService {
	return service.NewDashboardService(repo)
}

// ─────────────────────────────────────────────
// GetDashboardData
// ─────────────────────────────────────────────

func TestDashboardService_GetDashboardData_RepoError(t *testing.T) {
	repo := new(mocks.MockDashboardRepository)
	svc := newDashboardService(repo)

	repo.On("GetDashboardData", context.Background(), int64(1)).
		Return(nil, util.NewInternalError("db error"))

	result, err := svc.GetDashboardData(context.Background(), 1)

	require.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestDashboardService_GetDashboardData_Success_RoundsQuickInsights(t *testing.T) {
	repo := new(mocks.MockDashboardRepository)
	svc := newDashboardService(repo)

	curr := 12345.6789
	prev := 9876.5432
	total := 11111.1115
	repo.On("GetDashboardData", context.Background(), int64(1)).
		Return(&dto.DashboardDataDto{
			CurrentInvestmentValue:     &curr,
			PreviousDayInvestmentValue: &prev,
			TotalInvestedValue:         &total,
		}, nil)

	result, err := svc.GetDashboardData(context.Background(), 1)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 12345.68, *result.CurrentInvestmentValue)
	assert.Equal(t, 9876.54, *result.PreviousDayInvestmentValue)
	assert.Equal(t, 11111.11, *result.TotalInvestedValue)
	repo.AssertExpectations(t)
}

func TestDashboardService_GetDashboardData_Success_RoundsTopHoldings(t *testing.T) {
	repo := new(mocks.MockDashboardRepository)
	svc := newDashboardService(repo)

	repo.On("GetDashboardData", context.Background(), int64(1)).
		Return(&dto.DashboardDataDto{
			TopHoldings: []dto.HoldingResponseDto{
				{
					AssetName:       "HDFC Flexi Cap",
					Quantity:        10,
					AveragePrice:    123.456,
					CurrentPrice:    150.789,
					PrevDayPrice:    148.123,
					InvestedCapital: 1234.56,
					CurrentCapital:  0, // will be recalculated
				},
			},
		}, nil)

	result, err := svc.GetDashboardData(context.Background(), 1)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.TopHoldings, 1)
	h := result.TopHoldings[0]
	assert.Equal(t, 123.46, h.AveragePrice)
	assert.Equal(t, 150.79, h.CurrentPrice)
	assert.Equal(t, 148.12, h.PrevDayPrice)
	assert.Equal(t, 1234.56, h.InvestedCapital)
	// CurrentCapital = round2(CurrentPrice * Quantity) = round2(150.79 * 10)
	assert.Equal(t, 1507.9, h.CurrentCapital)
	// ReturnPercentage = round2((1507.9 - 1234.56) / 1234.56 * 100)
	assert.Equal(t, round2((1507.9-1234.56)/1234.56*100), h.ReturnPercentage)
	repo.AssertExpectations(t)
}

func TestDashboardService_GetDashboardData_Success_ZeroInvestedCapital(t *testing.T) {
	repo := new(mocks.MockDashboardRepository)
	svc := newDashboardService(repo)

	repo.On("GetDashboardData", context.Background(), int64(1)).
		Return(&dto.DashboardDataDto{
			TopHoldings: []dto.HoldingResponseDto{
				{
					AssetName:       "SBI Bluechip",
					Quantity:        5,
					AveragePrice:    0,
					CurrentPrice:    100,
					InvestedCapital: 0,
				},
			},
		}, nil)

	result, err := svc.GetDashboardData(context.Background(), 1)

	require.NoError(t, err)
	require.Len(t, result.TopHoldings, 1)
	assert.Equal(t, 0.0, result.TopHoldings[0].ReturnPercentage)
	repo.AssertExpectations(t)
}

func TestDashboardService_GetDashboardData_Success_RoundsTransactionPrice(t *testing.T) {
	repo := new(mocks.MockDashboardRepository)
	svc := newDashboardService(repo)

	desc := "buy"
	repo.On("GetDashboardData", context.Background(), int64(1)).
		Return(&dto.DashboardDataDto{
			RecentTransactions: []dto.TransactionResponseDto{
				{AssetName: "Infosys", TxnType: "buy", Quantity: 2, Price: 1500.555, Description: &desc, TxnDate: time.Now()},
				{AssetName: "TCS", TxnType: "sell", Quantity: 1, Price: 3200.111, TxnDate: time.Now()},
			},
		}, nil)

	result, err := svc.GetDashboardData(context.Background(), 1)

	require.NoError(t, err)
	require.Len(t, result.RecentTransactions, 2)
	assert.Equal(t, 1500.56, result.RecentTransactions[0].Price)
	assert.Equal(t, 3200.11, result.RecentTransactions[1].Price)
	// Quantity should remain unchanged
	assert.Equal(t, 2.0, result.RecentTransactions[0].Quantity)
	assert.Equal(t, 1.0, result.RecentTransactions[1].Quantity)
	repo.AssertExpectations(t)
}

func TestDashboardService_GetDashboardData_Success_EmptyHoldingsAndTransactions(t *testing.T) {
	repo := new(mocks.MockDashboardRepository)
	svc := newDashboardService(repo)

	zero := 0.0
	repo.On("GetDashboardData", context.Background(), int64(1)).
		Return(&dto.DashboardDataDto{
			CurrentInvestmentValue: &zero,
			TopHoldings:            []dto.HoldingResponseDto{},
			RecentTransactions:     []dto.TransactionResponseDto{},
		}, nil)

	result, err := svc.GetDashboardData(context.Background(), 1)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Empty(t, result.TopHoldings)
	assert.Empty(t, result.RecentTransactions)
	repo.AssertExpectations(t)
}

// round2 mirrors the helper in holding_service.go for assertion use.
func round2(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
}
