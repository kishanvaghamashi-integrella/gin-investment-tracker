package cron

import (
	"gin-investment-tracker/internal/cron/jobs"
	assetprice "gin-investment-tracker/internal/external-services/asset-details"
	repository "gin-investment-tracker/internal/repositories"
	"log"
	"log/slog"

	"github.com/robfig/cron/v3"
)

type CronJobs struct {
	assetRepo         repository.AssetRepositoryInterface
	priceDetailRepo   repository.PriceDetailRepositoryInterface
	assetPriceFetcher *assetprice.AssetPriceService
}

func NewCronJobs(assetRepo repository.AssetRepositoryInterface, priceDetailRepo repository.PriceDetailRepositoryInterface, assetPriceFetcher *assetprice.AssetPriceService) *CronJobs {
	return &CronJobs{assetRepo: assetRepo, priceDetailRepo: priceDetailRepo, assetPriceFetcher: assetPriceFetcher}
}

func (cj *CronJobs) Start() {
	c := cron.New(cron.WithSeconds())

	// Run at 12:00 AM everyday
	_, err := c.AddFunc("0 0 0 * * *", func() {
		slog.Info("Cron Job Started")
		jobs.FetchPriceDetailsJob(cj.assetRepo, cj.priceDetailRepo, cj.assetPriceFetcher)
		slog.Info("Cron Job Finished")
	})
	if err != nil {
		log.Fatal(err)
	}

	c.Start()
}
