package model

import "time"

type PriceDetail struct {
	ID        int64     `json:"id"`
	AssetID   int64     `json:"asset_id"`
	CurrPrice float64   `json:"curr_price"`
	PrevPrice float64   `json:"prev_price"`
	UpdatedAt time.Time `json:"updated_at"`
}
