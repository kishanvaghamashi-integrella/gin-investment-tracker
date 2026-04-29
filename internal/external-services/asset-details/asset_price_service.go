package assetprice

import (
	"fmt"
	"log/slog"
)

type AssetPriceService struct {
	stockFetcher AssetPriceFetcherInterface
	mfFetcher    AssetPriceFetcherInterface
}

func NewAssetPriceService(stockFetcher, mfFetcher AssetPriceFetcherInterface) *AssetPriceService {
	return &AssetPriceService{
		stockFetcher: stockFetcher,
		mfFetcher:    mfFetcher,
	}
}

func (s *AssetPriceService) FetchPrice(instrumentType, externalID string) (float64, float64, error) {
	slog.Info("External API called to fetch latest price")
	switch instrumentType {
	case "stock":
		return s.stockFetcher.FetchPrice(externalID)

	case "mutual_fund":
		return s.mfFetcher.FetchPrice(externalID)

	default:
		return 0, 0, fmt.Errorf("unsupported instrument type: %s", instrumentType)
	}
}
