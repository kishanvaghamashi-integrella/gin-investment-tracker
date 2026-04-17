package dto

type HoldingResponseDto struct {
	ID                      int64   `json:"id"`
	AssetExternalPlatformID string  `json:"external_platform_id" `
	AssetName               string  `json:"asset_name"`
	AssetInstrumentType     string  `json:"asset_instrument_type"`
	Quantity                float64 `json:"quantity"`
	AveragePrice            float64 `json:"average_price"`
	CurrentPrice            float64 `json:"current_price"`
	InvestedCapital         float64 `json:"invested_capital"`
	CurrentCapital          float64 `json:"current_capital"`
	ReturnPercentage        float64 `json:"return_percentage"`
}
