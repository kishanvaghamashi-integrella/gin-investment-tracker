package service

import (
	"context"
	"fmt"
	dto "gin-investment-tracker/internal/dtos"
	model "gin-investment-tracker/internal/models"
	repository "gin-investment-tracker/internal/repositories"
	"gin-investment-tracker/internal/util"
)

type UserAssetService struct {
	repo      repository.UserAssetRepositoryInterface
	userRepo  repository.UserRepositoryInterface
	assetRepo repository.AssetRepositoryInterface
}

func NewUserAssetService(
	repo repository.UserAssetRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
	assetRepo repository.AssetRepositoryInterface,
) *UserAssetService {
	return &UserAssetService{repo: repo, userRepo: userRepo, assetRepo: assetRepo}
}

func (s *UserAssetService) Create(ctx context.Context, userID int64, req *dto.CreateUserAssetRequest) (*model.UserAsset, error) {
	if err := s.ensureUserExists(ctx, userID); err != nil {
		return nil, err
	}

	exists, err := s.assetRepo.ExistsByID(ctx, req.AssetID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, util.NewNotFoundError(fmt.Sprintf("asset with id %d not found", req.AssetID))
	}

	exists, err = s.repo.IsUserAssetExits(ctx, userID, req.AssetID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, util.NewBadRequestError("This entry already exists")
	}

	userAsset := &model.UserAsset{
		UserID:  userID,
		AssetID: req.AssetID,
	}

	if err := s.repo.Create(ctx, userAsset); err != nil {
		return nil, err
	}

	return userAsset, nil
}

func (s *UserAssetService) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]model.UserAsset, error) {
	if err := s.ensureUserExists(ctx, userID); err != nil {
		return nil, err
	}

	return s.repo.GetByUserID(ctx, userID, limit, offset)
}

func (s *UserAssetService) Delete(ctx context.Context, userID, userAssetID int64) error {
	if err := s.ensureUserExists(ctx, userID); err != nil {
		return err
	}

	return s.repo.Delete(ctx, userAssetID, userID)
}

func (s *UserAssetService) ensureUserExists(ctx context.Context, userID int64) error {
	exists, err := s.userRepo.ExistsByID(ctx, userID)
	if err != nil {
		return err
	}
	if !exists {
		return util.NewNotFoundError(fmt.Sprintf("user with id %d not found", userID))
	}
	return nil
}
