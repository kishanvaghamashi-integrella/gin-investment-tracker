package jobs

import (
	"context"
	assetpricefetcher "gin-investment-tracker/internal/external-services/asset-details-fetcher"
	model "gin-investment-tracker/internal/models"
	repository "gin-investment-tracker/internal/repositories"
	"gin-investment-tracker/internal/util"
	"time"
)

func FetchPriceDetailsJob(assetRepo repository.AssetRepositoryInterface, priceDetailRepo repository.PriceDetailRepositoryInterface, assetPriceFetcher *assetpricefetcher.AssetPriceService) {
	log := util.Logger.With("job", "FetchPriceDetailsJob")
	logCtx := util.WithLogger(context.Background(), log)

	// 1. fetch all the assets by limit and offset
	limit := 50
	offset := 0

	for {
		ctx, cancel := context.WithTimeout(logCtx, 10*time.Second)
		assets, err := assetRepo.GetAll(ctx, limit, offset)
		cancel()
		offset += limit
		if err != nil {
			log.Errorw("Failed to fetch one batch of assets from DB", "limit", limit, "offset", offset-limit, "error", err)
			continue
		}
		if len(assets) == 0 {
			break
		}

		// 2. fetch price detail of each asset one by one using api
		var priceDetailList []model.PriceDetail
		for _, asset := range assets {
			if asset.ExternalPlatformID == nil {
				log.Warnw("skipping asset: missing external platform id", "asset_id", asset.ID, "instrument_type", asset.InstrumentType)
				continue
			}

			price, prevPrice, err := assetPriceFetcher.FetchAssetPrice(logCtx, asset.InstrumentType, *asset.ExternalPlatformID)
			if err != nil {
				log.Errorw("Failed to fetch price value for asset", "external_platform_id", asset.ExternalPlatformID, "error", err)
				continue
			}

			priceDetail := model.PriceDetail{
				AssetID:   asset.ID,
				CurrPrice: price,
				PrevPrice: prevPrice,
				UpdatedAt: time.Now(),
			}
			priceDetailList = append(priceDetailList, priceDetail)
			time.Sleep(10 * time.Microsecond)

		}

		// 3. add that price details into price_details table
		ctx2, cancel2 := context.WithTimeout(logCtx, 10*time.Second)
		err = priceDetailRepo.UpsertPriceDetails(ctx2, priceDetailList)
		cancel2()
		if err != nil {
			log.Errorw("Failed to add price value for into the DB", "limit", limit, "offset", offset-limit, "error", err)
			continue
		}
		log.Infow("Added price details successfully", "limit", limit, "offset", offset-limit)
	}
}
