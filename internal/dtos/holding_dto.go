package dto

type HoldingResponseDto struct {
	ID                  int64   `json:"id,omitempty"`
	AssetID             int64   `json:"asset_id,omitempty"`
	AssetName           string  `json:"asset_name,omitempty"`
	AssetInstrumentType string  `json:"asset_instrument_type,omitempty"`
	Quantity            float64 `json:"quantity,omitempty"`
	AveragePrice        float64 `json:"average_price,omitempty"`
	CurrentPrice        float64 `json:"current_price,omitempty"`
	PrevDayPrice        float64 `json:"prev_day_price,omitempty"`
	InvestedCapital     float64 `json:"invested_capital,omitempty"`
	CurrentCapital      float64 `json:"current_capital,omitempty"`
	ReturnPercentage    float64 `json:"return_percentage,omitempty"`
}
