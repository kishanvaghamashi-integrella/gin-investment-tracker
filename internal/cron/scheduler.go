package cron

import (
	"gin-investment-tracker/internal/cron/jobs"
	repository "gin-investment-tracker/internal/repositories"
	"log"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

type CronJobs struct {
	db              *pgxpool.Pool
	assetRepo       repository.AssetRepositoryInterface
	priceDetailRepo repository.PriceDetailRepositoryInterface
}

func NewCronJobs(assetRepo repository.AssetRepositoryInterface, priceDetailRepo repository.PriceDetailRepositoryInterface) *CronJobs {
	return &CronJobs{assetRepo: assetRepo, priceDetailRepo: priceDetailRepo}
}

func (cj *CronJobs) Start() {
	c := cron.New(cron.WithSeconds())

	// Run at 12:00 AM everyday
	_, err := c.AddFunc("0 56 11 * * *", func() {
		slog.Info("Cron Job Started")
		jobs.FetchPriceDetailsJob(cj.db, cj.assetRepo, cj.priceDetailRepo)
		slog.Info("Cron Job Finished")
	})
	if err != nil {
		log.Fatal(err)
	}

	c.Start()
}
