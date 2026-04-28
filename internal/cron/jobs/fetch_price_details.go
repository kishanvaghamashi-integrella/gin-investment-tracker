package jobs

import (
	"context"
	model "gin-investment-tracker/internal/models"
	repository "gin-investment-tracker/internal/repositories"
	service "gin-investment-tracker/internal/services"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func FetchPriceDetailsJob(db *pgxpool.Pool, assetRepo repository.AssetRepositoryInterface, priceDetailRepo repository.PriceDetailRepositoryInterface) {
	// 1. fetch all the assets by limit and offset
	limit := 50
	offset := 0

	for true {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		assets, err := assetRepo.GetAll(ctx, limit, offset)
		offset += limit
		if err != nil {
			slog.Error("Failed to fetch one batch of assets from DB", "limit", limit, "offset", offset-limit)
			continue
		}
		if len(assets) == 0 {
			break
		}

		// 2. fetch price detail of each asset one by one using api
		var priceDetailList []model.PriceDetail
		for _, asset := range assets {
			price, prevPrice, err := service.FetchLatestMfPrice(*asset.ExternalPlatformID)
			if err != nil {
				slog.Error("Failed to fetch price value for asset", "external platform ID", asset.ExternalPlatformID)
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
		ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel2()

		if err = priceDetailRepo.UpsertPriceDetails(ctx2, priceDetailList); err != nil {
			slog.Error("Failed to add price value for into the DB", "limit", limit, "offset", offset-limit, "error", err.Error())
			continue
		}
		slog.Info("Added price details successfully", "limit", limit, "offset", offset-limit)
	}
}
