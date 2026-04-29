package service

import (
	"context"
	"fmt"
	dto "gin-investment-tracker/internal/dtos"
	repository "gin-investment-tracker/internal/repositories"
	"gin-investment-tracker/internal/util"
	"math"
)

type HoldingService struct {
	repo     repository.HoldingRepositoryInterface
	userRepo repository.UserRepositoryInterface
}

func NewHoldingService(
	repo repository.HoldingRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
) *HoldingService {
	return &HoldingService{repo: repo, userRepo: userRepo}
}

func (s *HoldingService) GetAllByUserID(ctx context.Context, userID int64, limit, offset int) ([]dto.HoldingResponseDto, error) {
	if err := s.ensureUserExists(ctx, userID); err != nil {
		return nil, err
	}

	holdings, err := s.repo.GetAllByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	calculateProfitDetails(&holdings)
	return holdings, nil
}

func (s *HoldingService) ensureUserExists(ctx context.Context, userID int64) error {
	exists, err := s.userRepo.ExistsByID(ctx, userID)
	if err != nil {
		return err
	}
	if !exists {
		return util.NewNotFoundError(fmt.Sprintf("user with id %d not found", userID))
	}
	return nil
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func calculateProfitDetails(holdings *[]dto.HoldingResponseDto) {
	for i := range *holdings {
		h := &(*holdings)[i]
		h.Quantity = round2(h.Quantity)
		h.AveragePrice = round2(h.AveragePrice)
		h.CurrentPrice = round2(h.CurrentPrice)
		h.PrevDayPrice = round2(h.PrevDayPrice)
		h.InvestedCapital = round2(h.InvestedCapital)
		h.CurrentCapital = round2(h.CurrentPrice * h.Quantity)
		if h.InvestedCapital == 0 {
			h.ReturnPercentage = 0
		} else {
			h.ReturnPercentage = round2(((h.CurrentCapital - h.InvestedCapital) / h.InvestedCapital) * 100)
		}
	}
}
