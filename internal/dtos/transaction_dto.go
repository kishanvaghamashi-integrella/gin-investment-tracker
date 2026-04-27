package dto

import "time"

type CreateTransactionRequest struct {
	AssetID     int64     `json:"asset_id" binding:"required,gt=0"`
	TxnType     string    `json:"txn_type" binding:"required,txn_type"`
	Description *string   `json:"description"`
	Quantity    float64   `json:"quantity" binding:"required,gt=0"`
	Price       float64   `json:"price" binding:"required,gt=0"`
	TxnDate     time.Time `json:"txn_date" binding:"required"`
}

type UpdateTransactionRequest struct {
	TxnType     *string    `json:"txn_type" binding:"omitempty,txn_type"`
	Description *string    `json:"description"`
	Quantity    *float64   `json:"quantity" binding:"omitempty,gt=0"`
	Price       *float64   `json:"price" binding:"omitempty,gt=0"`
	TxnDate     *time.Time `json:"txn_date" binding:"omitempty"`
}

type ResponseTransactionDto struct {
	ID                  int64     `json:"id"`
	UserAssetID         int64     `json:"user_asset_id"`
	AssetName           string    `json:"asset_name"`
	AssetInstrumentType string    `json:"asset_instrument_type"`
	Description         *string   `json:"description"`
	TxnType             string    `json:"txn_type"`
	Quantity            float64   `json:"quantity"`
	Price               float64   `json:"price"`
	TxnDate             time.Time `json:"txn_date"`
}
