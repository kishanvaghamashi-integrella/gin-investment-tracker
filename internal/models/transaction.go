package model

import "time"

type Transaction struct {
	ID          int64     `json:"id" db:"id"`
	UserAssetID int64     `json:"user_asset_id" db:"user_asset_id"`
	TxnType     string    `json:"txn_type" db:"txn_type"`
	Quantity    float64   `json:"quantity" db:"quantity"`
	Price       float64   `json:"price" db:"price"`
	TxnDate     time.Time `json:"txn_date" db:"txn_date"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
