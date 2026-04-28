package service

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	dto "gin-investment-tracker/internal/dtos"
	model "gin-investment-tracker/internal/models"
	repository "gin-investment-tracker/internal/repositories"
	"gin-investment-tracker/internal/util"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// PriceFetcher abstracts external price lookups so the service can be tested
// deterministically without hitting live APIs.
type PriceFetcher interface {
	FetchPrice(instrumentType, externalPlatformID string) (float64, float64, error)
}

// RealPriceFetcher calls the live mfapi.in / Yahoo Finance endpoints.
type RealPriceFetcher struct{}

func (f *RealPriceFetcher) FetchPrice(instrumentType, externalPlatformID string) (float64, float64, error) {
	if instrumentType == "stock" {
		return FetchLatestStockPrice(externalPlatformID)
	}
	return FetchLatestMfPrice(externalPlatformID)
}

type HoldingService struct {
	repo         repository.HoldingRepositoryInterface
	userRepo     repository.UserRepositoryInterface
	priceFetcher PriceFetcher
}

func NewHoldingService(
	repo repository.HoldingRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
) *HoldingService {
	return &HoldingService{repo: repo, userRepo: userRepo, priceFetcher: &RealPriceFetcher{}}
}

// NewHoldingServiceWithFetcher is intended for tests that need a stubbed price fetcher.
func NewHoldingServiceWithFetcher(
	repo repository.HoldingRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
	fetcher PriceFetcher,
) *HoldingService {
	return &HoldingService{repo: repo, userRepo: userRepo, priceFetcher: fetcher}
}

func (s *HoldingService) GetAllByUserID(ctx context.Context, userID int64, limit, offset int) ([]dto.HoldingResponseDto, error) {
	if err := s.ensureUserExists(ctx, userID); err != nil {
		return nil, err
	}

	holdings, err := s.repo.GetAllByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	fetchLatestPriceAndCalculateHolding(&holdings, s.priceFetcher)
	return holdings, nil
}

func (s *HoldingService) ensureUserExists(ctx context.Context, userID int64) error {
	exists, err := s.userRepo.ExistsByID(ctx, userID)
	if err != nil {
		return err
	}
	if !exists {
		return util.NewNotFoundError(fmt.Sprintf("user with id %d not found", userID))
	}
	return nil
}

type priceResult struct {
	index int
	price float64
	err   error
}

func FetchLatestMfPrice(schemeCodeStr string) (float64, float64, error) {
	schemeCode, err := strconv.ParseInt(schemeCodeStr, 10, 64)
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

func FetchLatestStockPrice(schemeCode string) (float64, float64, error) {
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?interval=1d&range=1d", schemeCode)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create request for %s: %w", schemeCode, err)
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
		return 0, 0, fmt.Errorf("http request failed for scheme %s: %w", schemeCode, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("unexpected status %d for scheme %s", resp.StatusCode, schemeCode)
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
		return 0, 0, fmt.Errorf("failed to decode response for %s: %w", schemeCode, err)
	}

	if len(apiResp.Chart.Result) == 0 {
		return 0, 0, fmt.Errorf("no data returned for stock %s", schemeCode)
	}

	price := apiResp.Chart.Result[0].Meta.ChartPreviousClose
	return price, price, nil
}

func fetchLatestPriceAndCalculateHolding(holdings *[]dto.HoldingResponseDto, fetcher PriceFetcher) {
	resultCh := make(chan priceResult, len(*holdings))
	var wg sync.WaitGroup

	for i, holding := range *holdings {
		wg.Add(1)
		go func(i int, schemeCode, instrumentType string) {
			defer wg.Done()
			price, _, err := fetcher.FetchPrice(instrumentType, schemeCode)
			resultCh <- priceResult{index: i, price: price, err: err}
		}(i, holding.AssetExternalPlatformID, holding.AssetInstrumentType)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	for res := range resultCh {
		if res.err != nil {
			slog.Error(
				"failed to fetch current asset data",
				"assetName", (*holdings)[res.index].AssetName,
				"instrumentType", (*holdings)[res.index].AssetInstrumentType,
				"error", res.err,
			)
			res.price = (*holdings)[res.index].AveragePrice
			calculateCurrentProfit(res, holdings)
			continue
		}

		calculateCurrentProfit(res, holdings)
	}
}

func calculateCurrentProfit(priceData priceResult, holdings *[]dto.HoldingResponseDto) {
	holding := &(*holdings)[priceData.index]
	(*holding).CurrentPrice = priceData.price
	(*holding).CurrentCapital = priceData.price * (holding.Quantity)
	if holding.InvestedCapital == 0 {
		(*holding).ReturnPercentage = 0
	} else {

		(*holding).ReturnPercentage = ((holding.CurrentCapital - holding.InvestedCapital) / holding.InvestedCapital) * 100
	}
}
