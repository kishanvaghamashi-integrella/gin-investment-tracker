package repository

import (
	"context"

	model "gin-investment-tracker/internal/models"
)

type PriceDetailRepositoryInterface interface {
	UpsertPriceDetails(ctx context.Context, priceDetails []model.PriceDetail) error
}
