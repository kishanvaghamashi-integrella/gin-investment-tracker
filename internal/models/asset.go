package model

import "time"

type Asset struct {
	ID                 int64     `json:"id" db:"id"`
	Symbol             string    `json:"symbol" db:"symbol"`
	Name               string    `json:"name" db:"name"`
	InstrumentType     string    `json:"instrument_type" db:"instrument_type"`
	ISIN               string    `json:"isin" db:"isin"`
	Exchange           string    `json:"exchange" db:"exchange"`
	Currency           string    `json:"currency" db:"currency"`
	ExternalPlatformID string    `json:"external_platform_id" db:"external_platform_id"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}
