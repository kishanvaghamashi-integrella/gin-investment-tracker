package model

import "time"

type UserAsset struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	AssetID   int64     `json:"asset_id" db:"asset_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
