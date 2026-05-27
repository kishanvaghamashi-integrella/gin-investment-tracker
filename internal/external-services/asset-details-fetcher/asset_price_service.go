package assetpricefetcher

import (
	"context"
	"fmt"
	"gin-investment-tracker/internal/util"
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

func (s *AssetPriceService) FetchAssetPrice(ctx context.Context, instrumentType, externalID string) (float64, float64, error) {
	util.FromContext(ctx).With("service", "AssetPriceService.FetchAssetPrice").Infow("External API called to fetch latest price", "externalID", externalID)
	switch instrumentType {
	case "stock":
		return s.stockFetcher.FetchPrice(externalID)

	case "mutual_fund":
		return s.mfFetcher.FetchPrice(externalID)

	default:
		return 0, 0, fmt.Errorf("unsupported instrument type: %s", instrumentType)
	}
}
