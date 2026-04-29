package dto

type HoldingResponseDto struct {
	ID                  int64   `json:"id"`
	AssetID             int64   `json:"asset_id"`
	AssetName           string  `json:"asset_name"`
	AssetInstrumentType string  `json:"asset_instrument_type"`
	Quantity            float64 `json:"quantity"`
	AveragePrice        float64 `json:"average_price"`
	CurrentPrice        float64 `json:"current_price"`
	PrevDayPrice        float64 `json:"prev_day_price"`
	InvestedCapital     float64 `json:"invested_capital"`
	CurrentCapital      float64 `json:"current_capital"`
	ReturnPercentage    float64 `json:"return_percentage"`
}
