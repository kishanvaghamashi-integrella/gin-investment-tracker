package dto

import "time"

type CreateTransactionRequest struct {
	AssetID     int64     `json:"asset_id,omitempty" binding:"required,gt=0"`
	TxnType     string    `json:"txn_type,omitempty" binding:"required,txn_type"`
	Description *string   `json:"description,omitempty"`
	Quantity    float64   `json:"quantity,omitempty" binding:"required,gt=0"`
	Price       float64   `json:"price,omitempty" binding:"required,gt=0"`
	TxnDate     time.Time `json:"txn_date,omitempty" binding:"required"`
}

type UpdateTransactionRequest struct {
	TxnType     *string    `json:"txn_type,omitempty" binding:"omitempty,txn_type"`
	Description *string    `json:"description,omitempty"`
	Quantity    *float64   `json:"quantity,omitempty" binding:"omitempty,gt=0"`
	Price       *float64   `json:"price,omitempty" binding:"omitempty,gt=0"`
	TxnDate     *time.Time `json:"txn_date,omitempty" binding:"omitempty"`
}

type TransactionResponseDto struct {
	ID                  int64     `json:"id,omitempty"`
	UserAssetID         int64     `json:"user_asset_id,omitempty"`
	AssetName           string    `json:"asset_name,omitempty"`
	AssetInstrumentType string    `json:"asset_instrument_type,omitempty"`
	Description         *string   `json:"description,omitempty"`
	TxnType             string    `json:"txn_type,omitempty"`
	Quantity            float64   `json:"quantity,omitempty"`
	Price               float64   `json:"price,omitempty"`
	TxnDate             time.Time `json:"txn_date,omitempty"`
}
