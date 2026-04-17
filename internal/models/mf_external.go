package model

type MfNavData struct {
	Date string `json:"date"`
	Nav  string `json:"nav"`
}

type MfNavApiResponse struct {
	Data   []MfNavData `json:"data"`
	Status string      `json:"status"`
}
