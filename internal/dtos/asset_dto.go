package dto

type CreateAssetRequest struct {
	Symbol             string  `json:"symbol" binding:"required,min=1,max=100"`
	Name               string  `json:"name" binding:"required,min=1,max=200"`
	InstrumentType     string  `json:"instrument_type" binding:"required,instrument_type"`
	ISIN               string  `json:"isin" binding:"required,min=10,max=50"`
	Exchange           string  `json:"exchange" binding:"required,min=2,max=100"`
	Currency           string  `json:"currency" binding:"omitempty,min=2,max=10"`
	ExternalPlatformID *string `json:"external_platform_id" binding:"omitempty,min=2,max=100"`
}

type UpdateAssetRequest struct {
	Symbol             *string `json:"symbol" binding:"omitempty,min=1,max=100"`
	Name               *string `json:"name" binding:"omitempty,min=1,max=200"`
	InstrumentType     *string `json:"instrument_type" binding:"omitempty,instrument_type"`
	ISIN               *string `json:"isin" binding:"omitempty,min=10,max=50"`
	Exchange           *string `json:"exchange" binding:"omitempty,min=2,max=100"`
	Currency           *string `json:"currency" binding:"omitempty,min=2,max=10"`
	ExternalPlatformID *string `json:"external_platform_id" binding:"omitempty,min=2,max=100"`
}
