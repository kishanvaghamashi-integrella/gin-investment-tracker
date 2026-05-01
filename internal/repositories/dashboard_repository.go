package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	dto "gin-investment-tracker/internal/dtos"
	"gin-investment-tracker/internal/util"
	"log/slog"
)

type DashboardRepository struct {
	db *pgxpool.Pool
}

func NewDashboardRepository(db *pgxpool.Pool) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) GetDashboardData(ctx context.Context, userID int64) (*dto.DashboardDataDto, error) {
	var dashboardDataDto dto.DashboardDataDto

	// Quick insights
	query := `
		SELECT 
			COALESCE(SUM(h.total_quantity * COALESCE(pd.curr_price, 0)), 0) as current_investment_value,
			COALESCE(SUM(h.total_quantity * COALESCE(pd.prev_price, 0)), 0) as previous_day_investment_value,
			COALESCE(SUM(h.total_invested), 0) as total_invested_value
		FROM holdings h
		INNER JOIN user_assets ua ON h.user_asset_id = ua.id
		LEFT JOIN price_details pd ON ua.asset_id = pd.asset_id
		WHERE ua.user_id = $1
	`
	var curr, prev, total float64
	if err := r.db.QueryRow(ctx, query, userID).Scan(
		&curr,
		&prev,
		&total,
	); err != nil {
		slog.Error("failed to fetch dashboard data quick insights of portfolio", "error", err.Error())
		return nil, util.NewInternalError("failed to fetch dashboard data")
	}
	dashboardDataDto.CurrentInvestmentValue = &curr
	dashboardDataDto.PreviousDayInvestmentValue = &prev
	dashboardDataDto.TotalInvestedValue = &total

	// Top Holdings
	query = `
		SELECT 
			a.name,
			a.instrument_type,
			h.total_quantity,
			h.average_price,
			pd.curr_price,
			h.total_invested,
			(h.total_quantity * pd.curr_price) as current_value
		FROM holdings h
		INNER JOIN user_assets ua ON h.user_asset_id = ua.id
		INNER JOIN assets a on ua.asset_id = a.id
		INNER JOIN price_details pd ON ua.asset_id = pd.asset_id
		WHERE ua.user_id = $1
		ORDER BY current_value DESC
		LIMIT 5;
	`
	rows1, err := r.db.Query(ctx, query, userID)
	if err != nil {
		slog.Error("failed to fetch top holdings data for dashboard", "error", err.Error())
		return nil, util.NewInternalError("failed to fetch dashboard data")
	}

	var holdings []dto.HoldingResponseDto
	for rows1.Next() {
		var h dto.HoldingResponseDto
		if err := rows1.Scan(
			&h.AssetName,
			&h.AssetInstrumentType,
			&h.Quantity,
			&h.AveragePrice,
			&h.CurrentPrice,
			&h.InvestedCapital,
			&h.CurrentCapital,
		); err != nil {
			rows1.Close()
			slog.Error("failed to scan holding row", "error", err.Error())
			return nil, util.NewInternalError("failed to fetch dashboard data")
		}
		holdings = append(holdings, h)
	}
	rows1.Close()

	if err := rows1.Err(); err != nil {
		slog.Error("failed to iterate holding rows", "error", err.Error())
		return nil, util.NewInternalError("failed to fetch dashboard data")
	}
	dashboardDataDto.TopHoldings = holdings

	// Top 5 transaction
	query = `
		SELECT  
			a.name, 
			a.instrument_type,
			description, 
			quantity, 
			price, 
			txn_date 
		FROM transactions t
		INNER JOIN user_assets ua ON t.user_asset_id = ua.id
		INNER JOIN assets a on ua.asset_id = a.id
		WHERE ua.user_id = $1
		ORDER BY txn_date DESC
		LIMIT 5
	`
	rows2, err := r.db.Query(ctx, query, userID)
	if err != nil {
		slog.Error("failed to fetch recent transactions data for dashboard", "error", err.Error())
		return nil, util.NewInternalError("failed to fetch dashboard data")
	}

	var transactions []dto.TransactionResponseDto
	for rows2.Next() {
		var t dto.TransactionResponseDto
		if err := rows2.Scan(
			&t.AssetName,
			&t.AssetInstrumentType,
			&t.Description,
			&t.Quantity,
			&t.Price,
			&t.TxnDate,
		); err != nil {
			rows2.Close()
			slog.Error("failed to scan transaction row", "error", err.Error())
			return nil, util.NewInternalError("failed to fetch dashboard data")
		}
		transactions = append(transactions, t)
	}
	rows2.Close()

	if err := rows2.Err(); err != nil {
		slog.Error("failed to iterate transaction rows", "error", err.Error())
		return nil, util.NewInternalError("failed to fetch dashboard data")
	}
	dashboardDataDto.RecentTransactions = transactions

	return &dashboardDataDto, nil
}
