package assetprice

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	model "gin-investment-tracker/internal/models"
	"io"
	"net/http"
	"time"
)

type YahooStockFetcher struct{}

func NewYahooStockFetcher() *YahooStockFetcher {
	return &YahooStockFetcher{}
}

func (y *YahooStockFetcher) FetchPrice(externalID string) (float64, float64, error) {
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?interval=1d&range=1d", externalID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create request for %s: %w", externalID, err)
	}

	// Mimicking a real browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://finance.yahoo.com/")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		return 0, 0, fmt.Errorf("http request failed for scheme %s: %w", externalID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("unexpected status %d for scheme %s", resp.StatusCode, externalID)
	}

	// Handle gzip encoding since we declared Accept-Encoding above
	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzReader.Close()
		reader = gzReader
	}

	var apiResp model.StockChartResult
	if err := json.NewDecoder(reader).Decode(&apiResp); err != nil {
		return 0, 0, fmt.Errorf("failed to decode response for %s: %w", externalID, err)
	}

	if len(apiResp.Chart.Result) == 0 {
		return 0, 0, fmt.Errorf("no data returned for stock %s", externalID)
	}

	price := apiResp.Chart.Result[0].Meta.ChartPreviousClose
	return price, price, nil
}
