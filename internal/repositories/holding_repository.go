package repository

import (
	"context"
	"fmt"
	dto "gin-investment-tracker/internal/dtos"
	"gin-investment-tracker/internal/util"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HoldingRepository struct {
	db *pgxpool.Pool
}

func NewHoldingRepository(db *pgxpool.Pool) *HoldingRepository {
	return &HoldingRepository{db: db}
}

func (r *HoldingRepository) GetAllByUserID(ctx context.Context, userID int64, limit, offset int, sortByQuery, assetNameQuery string) ([]dto.HoldingResponseDto, error) {
	switch strings.ToLower(sortByQuery) {
	case "id":
		sortByQuery = "h.id"
	case "name":
		sortByQuery = "a.name"
	case "profit":
		sortByQuery = "h.total_quantity * (COALESCE(pd.curr_price, 0) - h.average_price) DESC"
	default:
		sortByQuery = "h.id"
	}

	query := fmt.Sprintf(`
		SELECT h.id, a.id, a.name, a.instrument_type, h.total_quantity, h.average_price, COALESCE(pd.curr_price, 0), COALESCE(pd.prev_price, 0), h.total_invested
		FROM holdings h
		INNER JOIN user_assets ua ON h.user_asset_id = ua.id
		INNER JOIN assets a ON ua.asset_id = a.id
		LEFT JOIN price_details pd ON a.id = pd.asset_id
		WHERE ua.user_id = $1 AND LOWER(a.name) LIKE $4
		ORDER BY %s
		LIMIT $2 OFFSET $3
	`, sortByQuery)

	rows, err := r.db.Query(ctx, query, userID, limit, offset, "%"+assetNameQuery+"%")
	if err != nil {
		slog.Error("failed to list holdings", "error", err.Error())
		return nil, util.NewInternalError("failed to list holdings")
	}
	defer rows.Close()

	var holdings []dto.HoldingResponseDto
	for rows.Next() {
		var h dto.HoldingResponseDto
		if err := rows.Scan(&h.ID, &h.AssetID, &h.AssetName, &h.AssetInstrumentType, &h.Quantity, &h.AveragePrice, &h.CurrentPrice, &h.PrevDayPrice, &h.InvestedCapital); err != nil {
			slog.Error("failed to scan holding row", "error", err.Error())
			return nil, util.NewInternalError("failed to list holdings")
		}
		holdings = append(holdings, h)
	}

	if err := rows.Err(); err != nil {
		slog.Error("failed to iterate holding rows", "error", err.Error())
		return nil, util.NewInternalError("failed to list holdings")
	}

	return holdings, nil
}
