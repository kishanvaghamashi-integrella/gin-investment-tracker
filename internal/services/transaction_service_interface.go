package service

import (
	"context"

	dto "gin-investment-tracker/internal/dtos"
	model "gin-investment-tracker/internal/models"
)

type TransactionServiceInterface interface {
	Create(ctx context.Context, req *dto.CreateTransactionRequest, userId int64) (*model.Transaction, error)
	GetAllByUserID(ctx context.Context, userID int64, limit, offset int) ([]dto.ResponseTransactionDto, error)
	Update(ctx context.Context, id int64, req *dto.UpdateTransactionRequest) error
	Delete(ctx context.Context, id int64) error
}
