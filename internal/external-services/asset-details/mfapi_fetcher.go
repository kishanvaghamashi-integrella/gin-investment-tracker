package assetprice

import (
	"encoding/json"
	"fmt"
	model "gin-investment-tracker/internal/models"
	"net/http"
	"strconv"
	"time"
)

type MfapiFetcher struct{}

func NewMfapiFetcher() *MfapiFetcher {
	return &MfapiFetcher{}
}

func (m *MfapiFetcher) FetchPrice(externalID string) (float64, float64, error) {
	schemeCode, err := strconv.ParseInt(externalID, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert scheme code from string to integer")
	}

	currDate := time.Now()
	endDateStr := currDate.Format("2001-12-02")
	startDate := currDate.AddDate(0, 0, -7)
	startDateStr := startDate.Format("2001-12-02")

	url := fmt.Sprintf("https://api.mfapi.in/mf/%d?startDate=%s&endDate=%s", schemeCode, startDateStr, endDateStr)

	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, fmt.Errorf("http request failed for scheme %d: %w", schemeCode, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("unexpected status %d for scheme %d", resp.StatusCode, schemeCode)
	}

	var apiResp model.MfNavApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return 0, 0, fmt.Errorf("failed to decode response for scheme %d: %w", schemeCode, err)
	}

	if apiResp.Status != "SUCCESS" {
		return 0, 0, fmt.Errorf("api returned non-success status for scheme %d", schemeCode)
	}

	if len(apiResp.Data) == 0 {
		return 0, 0, fmt.Errorf("no NAV data returned for scheme %d", schemeCode)
	}

	price, err := strconv.ParseFloat(apiResp.Data[0].Nav, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse NAV value '%s': %w", apiResp.Data[0].Nav, err)
	}
	prevDayPrice, err := strconv.ParseFloat(apiResp.Data[1].Nav, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse previous day NAV value '%s': %w", apiResp.Data[0].Nav, err)
	}

	return price, prevDayPrice, nil
}
