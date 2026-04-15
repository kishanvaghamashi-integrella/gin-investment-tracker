package service

import (
	"context"
	dto "gin-investment-tracker/internal/dtos"
	model "gin-investment-tracker/internal/models"
	repository "gin-investment-tracker/internal/repositories"
)

type AssetService struct {
	repo repository.AssetRepositoryInterface
}

func NewAssetService(repo repository.AssetRepositoryInterface) *AssetService {
	return &AssetService{repo: repo}
}

func (s *AssetService) Create(ctx context.Context, req *dto.CreateAssetRequest) (*model.Asset, error) {
	asset := &model.Asset{
		Symbol:             req.Symbol,
		Name:               req.Name,
		InstrumentType:     req.InstrumentType,
		ISIN:               req.ISIN,
		Exchange:           req.Exchange,
		Currency:           req.Currency,
		ExternalPlatformID: req.ExternalPlatformID,
	}

	if asset.Currency == "" {
		asset.Currency = "INR"
	}

	if err := s.repo.Create(ctx, asset); err != nil {
		return nil, err
	}

	return asset, nil
}

func (s *AssetService) GetByID(ctx context.Context, id int64) (*model.Asset, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *AssetService) GetAll(ctx context.Context, limit, offset int) ([]model.Asset, error) {
	return s.repo.GetAll(ctx, limit, offset)
}

func (s *AssetService) Update(ctx context.Context, id int64, req *dto.UpdateAssetRequest) error {
	asset, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if req.Symbol != nil {
		asset.Symbol = *req.Symbol
	}
	if req.Name != nil {
		asset.Name = *req.Name
	}
	if req.InstrumentType != nil {
		asset.InstrumentType = *req.InstrumentType
	}
	if req.ISIN != nil {
		asset.ISIN = *req.ISIN
	}
	if req.Exchange != nil {
		asset.Exchange = *req.Exchange
	}
	if req.Currency != nil {
		asset.Currency = *req.Currency
	}
	if req.ExternalPlatformID != nil {
		asset.ExternalPlatformID = *req.ExternalPlatformID
	}

	return s.repo.Update(ctx, asset)
}

func (s *AssetService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
