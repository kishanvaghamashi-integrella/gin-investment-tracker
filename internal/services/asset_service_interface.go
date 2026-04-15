package service

import (
	"context"
	dto "gin-investment-tracker/internal/dtos"
	model "gin-investment-tracker/internal/models"
)

type AssetServiceInterface interface {
	Create(ctx context.Context, req *dto.CreateAssetRequest) (*model.Asset, error)
	GetByID(ctx context.Context, id int64) (*model.Asset, error)
	GetAll(ctx context.Context, limit, offset int) ([]model.Asset, error)
	Update(ctx context.Context, id int64, req *dto.UpdateAssetRequest) error
	Delete(ctx context.Context, id int64) error
}
