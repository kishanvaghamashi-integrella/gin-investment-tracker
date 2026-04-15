package repository

import (
	"context"
	"errors"
	"fmt"
	model "gin-investment-tracker/internal/models"
	"gin-investment-tracker/internal/util"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AssetRepository struct {
	db *pgxpool.Pool
}

func NewAssetRepository(db *pgxpool.Pool) *AssetRepository {
	return &AssetRepository{db: db}
}

func (r *AssetRepository) Create(ctx context.Context, asset *model.Asset) error {
	query := `
		INSERT INTO assets (symbol, name, instrument_type, isin, exchange, currency, external_platform_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		asset.Symbol,
		asset.Name,
		asset.InstrumentType,
		asset.ISIN,
		asset.Exchange,
		asset.Currency,
		asset.ExternalPlatformID,
	).Scan(&asset.ID, &asset.CreatedAt)

	if err != nil {
		slog.Error("failed to create asset", "error", err.Error())
		return util.NewInternalError("failed to create asset")
	}

	return nil
}

func (r *AssetRepository) GetByID(ctx context.Context, id int64) (*model.Asset, error) {
	query := `
		SELECT id, symbol, name, instrument_type, isin, exchange, currency, external_platform_id, created_at
		FROM assets
		WHERE id = $1
	`

	var asset model.Asset
	err := r.db.QueryRow(ctx, query, id).Scan(
		&asset.ID,
		&asset.Symbol,
		&asset.Name,
		&asset.InstrumentType,
		&asset.ISIN,
		&asset.Exchange,
		&asset.Currency,
		&asset.ExternalPlatformID,
		&asset.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, util.NewNotFoundError(fmt.Sprintf("asset with id %d not found", id))
		}
		slog.Error("failed to get asset", "error", err.Error())
		return nil, util.NewInternalError("failed to get asset")
	}

	return &asset, nil
}

func (r *AssetRepository) GetAll(ctx context.Context, limit, offset int) ([]model.Asset, error) {
	query := `
		SELECT id, symbol, name, instrument_type, isin, exchange, currency, external_platform_id, created_at
		FROM assets
		ORDER BY id
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		slog.Error("failed to list assets", "error", err.Error())
		return nil, util.NewInternalError("failed to list assets")
	}
	defer rows.Close()

	var assets []model.Asset
	for rows.Next() {
		var asset model.Asset
		if err := rows.Scan(
			&asset.ID,
			&asset.Symbol,
			&asset.Name,
			&asset.InstrumentType,
			&asset.ISIN,
			&asset.Exchange,
			&asset.Currency,
			&asset.ExternalPlatformID,
			&asset.CreatedAt,
		); err != nil {
			slog.Error("failed to scan asset row", "error", err.Error())
			return nil, util.NewInternalError("failed to list assets")
		}
		assets = append(assets, asset)
	}

	if err := rows.Err(); err != nil {
		slog.Error("failed while iterating asset rows", "error", err.Error())
		return nil, util.NewInternalError("failed to list assets")
	}

	return assets, nil
}

func (r *AssetRepository) Update(ctx context.Context, asset *model.Asset) error {
	query := `
		UPDATE assets
		SET symbol = $2, name = $3, instrument_type = $4, isin = $5, exchange = $6, currency = $7, external_platform_id = $8
		WHERE id = $1
	`

	res, err := r.db.Exec(
		ctx,
		query,
		asset.ID,
		asset.Symbol,
		asset.Name,
		asset.InstrumentType,
		asset.ISIN,
		asset.Exchange,
		asset.Currency,
		asset.ExternalPlatformID,
	)

	if err != nil {
		slog.Error("failed to update asset", "error", err.Error())
		return util.NewInternalError("failed to update asset")
	}

	if res.RowsAffected() == 0 {
		return util.NewNotFoundError(fmt.Sprintf("asset with id %d not found", asset.ID))
	}

	return nil
}

func (r *AssetRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM assets WHERE id = $1`

	res, err := r.db.Exec(ctx, query, id)
	if err != nil {
		slog.Error("failed to delete asset", "error", err.Error())
		return util.NewInternalError("failed to delete asset")
	}

	if res.RowsAffected() == 0 {
		return util.NewNotFoundError(fmt.Sprintf("asset with id %d not found", id))
	}

	return nil
}

func (r *AssetRepository) ExistsByID(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM assets WHERE id = $1)`

	var exists bool
	if err := r.db.QueryRow(ctx, query, id).Scan(&exists); err != nil {
		slog.Error("failed to check asset existence", "error", err.Error())
		return false, util.NewInternalError("failed to check asset existence")
	}

	return exists, nil
}
