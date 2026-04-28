package repository

import (
	"context"
	"time"

	model "gin-investment-tracker/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PriceDetailRepository struct {
	db *pgxpool.Pool
}

func NewPriceDetailRepository(db *pgxpool.Pool) *PriceDetailRepository {
	return &PriceDetailRepository{db: db}
}

func (r *PriceDetailRepository) UpsertPriceDetails(ctx context.Context, priceDetails []model.PriceDetail) error {
	assetIDs := make([]int64, len(priceDetails))
	currPrices := make([]float64, len(priceDetails))
	prevPrices := make([]float64, len(priceDetails))
	updatedAts := make([]time.Time, len(priceDetails))

	for i, pd := range priceDetails {
		assetIDs[i] = pd.AssetID
		currPrices[i] = pd.CurrPrice
		prevPrices[i] = pd.PrevPrice
		updatedAts[i] = pd.UpdatedAt
	}

	query := `
		INSERT INTO price_details (asset_id, curr_price, prev_price, updated_at)
		SELECT * FROM unnest($1::bigint[], $2::numeric[], $3::numeric[], $4::timestamptz[])
		ON CONFLICT (asset_id) DO UPDATE SET
			curr_price = EXCLUDED.curr_price,
			prev_price = EXCLUDED.prev_price,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.Exec(ctx, query, assetIDs, currPrices, prevPrices, updatedAts)
	return err
}
