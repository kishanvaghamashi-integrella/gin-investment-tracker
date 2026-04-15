package service

import (
	"context"

	dto "gin-investment-tracker/internal/dtos"
	model "gin-investment-tracker/internal/models"
)

type UserAssetServiceInterface interface {
	Create(ctx context.Context, userID int64, req *dto.CreateUserAssetRequest) (*model.UserAsset, error)
	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]model.UserAsset, error)
	Delete(ctx context.Context, userID, userAssetID int64) error
}
