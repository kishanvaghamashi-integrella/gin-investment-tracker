package dto

type CreateAssetRequest struct {
	Symbol             string  `json:"symbol,omitempty" binding:"required,min=1,max=100"`
	Name               string  `json:"name,omitempty" binding:"required,min=1,max=200"`
	InstrumentType     string  `json:"instrument_type,omitempty" binding:"required,instrument_type"`
	AMC                *string `json:"amc,omitempty"`
	ISIN               string  `json:"isin,omitempty" binding:"required,min=10,max=50"`
	Exchange           string  `json:"exchange,omitempty" binding:"required,min=2,max=100"`
	Currency           string  `json:"currency,omitempty" binding:"omitempty,min=2,max=10"`
	ExternalPlatformID *string `json:"external_platform_id,omitempty" binding:"omitempty,min=2,max=100"`
}

type UpdateAssetRequest struct {
	Symbol             *string `json:"symbol,omitempty" binding:"omitempty,min=1,max=100"`
	Name               *string `json:"name,omitempty" binding:"omitempty,min=1,max=200"`
	InstrumentType     *string `json:"instrument_type,omitempty" binding:"omitempty,instrument_type"`
	AMC                *string `json:"amc,omitempty"`
	ISIN               *string `json:"isin,omitempty" binding:"omitempty,min=10,max=50"`
	Exchange           *string `json:"exchange,omitempty" binding:"omitempty,min=2,max=100"`
	Currency           *string `json:"currency,omitempty" binding:"omitempty,min=2,max=10"`
	ExternalPlatformID *string `json:"external_platform_id,omitempty" binding:"omitempty,min=2,max=100"`
}
