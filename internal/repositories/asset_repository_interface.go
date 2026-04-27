package repository

import (
	"context"
	model "gin-investment-tracker/internal/models"
)

type AssetRepositoryInterface interface {
	Create(ctx context.Context, asset *model.Asset) error
	GetByID(ctx context.Context, id int64) (*model.Asset, error)
	GetByISIN(ctx context.Context, isin string) (*model.Asset, error)
	GetAll(ctx context.Context, limit, offset int) ([]model.Asset, error)
	Update(ctx context.Context, asset *model.Asset) error
	Delete(ctx context.Context, id int64) error
	ExistsByID(ctx context.Context, id int64) (bool, error)
}
