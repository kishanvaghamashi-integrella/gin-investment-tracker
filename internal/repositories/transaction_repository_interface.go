package repository

import (
	"context"
	dto "gin-investment-tracker/internal/dtos"
	model "gin-investment-tracker/internal/models"
)

type TransactionRepositoryInterface interface {
	Create(ctx context.Context, txn *model.Transaction, holding *model.Holding, isUpdate bool) error
	GetAllByUserID(ctx context.Context, userID int64, limit, offset int) ([]dto.ResponseTransactionDto, error)
	GetHoldingsByUserAssetID(ctx context.Context, userAssetID int64) (*model.Holding, error)
	GetByID(ctx context.Context, id int64) (*model.Transaction, error)
	Update(ctx context.Context, txn *model.Transaction) error
	Delete(ctx context.Context, id int64) error
}
