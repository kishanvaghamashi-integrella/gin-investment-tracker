package model

import "time"

type Holding struct {
	ID            int64     `json:"id" db:"id"`
	UserAssetID   int64     `json:"user_asset_id" db:"user_asset_id"`
	TotalQuantity float64   `json:"total_quantity" db:"total_quantity"`
	AveragePrice  float64   `json:"average_price" db:"average_price"`
	TotalInvested float64   `json:"total_invested" db:"total_invested"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
