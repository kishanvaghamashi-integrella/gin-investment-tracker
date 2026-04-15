package repository

import (
	"context"
	model "gin-investment-tracker/internal/models"
)

type UserAssetRepositoryInterface interface {
	Create(ctx context.Context, userAsset *model.UserAsset) error
	GetIdByUserIdAssetId(ctx context.Context, userID, assetID int64) (*int64, error)
	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]model.UserAsset, error)
	Delete(ctx context.Context, id, userID int64) error
	IsUserAssetExits(ctx context.Context, userID int64, assetID int64) (bool, error)
	ExistsByID(ctx context.Context, id int64) (bool, error)
}
