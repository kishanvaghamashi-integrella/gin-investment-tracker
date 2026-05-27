package cron

import (
	"gin-investment-tracker/internal/cron/jobs"
	assetprice "gin-investment-tracker/internal/external-services/asset-details"
	repository "gin-investment-tracker/internal/repositories"
	"gin-investment-tracker/internal/util"

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

func (cj *CronJobs) Start() error {
	c := cron.New(cron.WithSeconds())

	// Run at 12:00 AM everyday
	_, err := c.AddFunc("0 0 0 * * *", func() {
		util.Logger.Infow("cron job started")
		jobs.FetchPriceDetailsJob(cj.assetRepo, cj.priceDetailRepo, cj.assetPriceFetcher)
		util.Logger.Infow("cron job finished")
	})
	if err != nil {
		return err
	}

	c.Start()
	return nil
}
