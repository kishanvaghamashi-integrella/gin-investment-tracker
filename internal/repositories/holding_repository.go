package repository

import (
	"context"
	dto "gin-investment-tracker/internal/dtos"
	"gin-investment-tracker/internal/util"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HoldingRepository struct {
	db *pgxpool.Pool
}

func NewHoldingRepository(db *pgxpool.Pool) *HoldingRepository {
	return &HoldingRepository{db: db}
}

func (r *HoldingRepository) GetAllByUserID(ctx context.Context, userID int64, limit, offset int) ([]dto.HoldingResponseDto, error) {
	query := `
		SELECT h.id, a.external_platform_id, a.name, a.instrument_type, h.total_quantity, h.average_price, h.total_invested
		FROM holdings h
		INNER JOIN user_assets ua ON h.user_asset_id = ua.id
		INNER JOIN assets a ON a.id = ua.asset_id
		WHERE ua.user_id = $1
		ORDER BY h.id
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		slog.Error("failed to list holdings", "error", err.Error())
		return nil, util.NewInternalError("failed to list holdings")
	}
	defer rows.Close()

	var holdings []dto.HoldingResponseDto
	for rows.Next() {
		var h dto.HoldingResponseDto
		if err := rows.Scan(&h.ID, &h.AssetExternalPlatformID, &h.AssetName, &h.AssetInstrumentType, &h.Quantity, &h.AveragePrice, &h.InvestedCapital); err != nil {
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
