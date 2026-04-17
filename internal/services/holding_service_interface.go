package service

import (
	"context"

	dto "gin-investment-tracker/internal/dtos"
)

type HoldingServiceInterface interface {
	GetAllByUserID(ctx context.Context, userID int64, limit, offset int) ([]dto.HoldingResponseDto, error)
}
