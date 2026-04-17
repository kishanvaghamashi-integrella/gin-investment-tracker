package repository

import (
	"context"
	dto "gin-investment-tracker/internal/dtos"
)

type HoldingRepositoryInterface interface {
	GetAllByUserID(ctx context.Context, userID int64, limit, offset int) ([]dto.HoldingResponseDto, error)
}
