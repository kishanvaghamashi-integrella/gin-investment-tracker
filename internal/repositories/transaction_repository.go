package repository

import (
	"context"
	"fmt"
	dto "gin-investment-tracker/internal/dtos"
	model "gin-investment-tracker/internal/models"
	"gin-investment-tracker/internal/util"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, txn *model.Transaction, holding *model.Holding, isUpdate bool) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		slog.Error("failed to begin transaction", "error", err.Error())
		return util.NewInternalError("failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	txnQuery := `
		INSERT INTO transactions (user_asset_id, txn_type, quantity, price, txn_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	err = tx.QueryRow(ctx, txnQuery, txn.UserAssetID, txn.TxnType, txn.Quantity, txn.Price, txn.TxnDate).
		Scan(&txn.ID, &txn.CreatedAt)
	if err != nil {
		slog.Error("failed to insert transaction", "error", err.Error())
		return util.NewInternalError("failed to create transaction")
	}

	if isUpdate {
		holdingQuery := `
			UPDATE holdings
			SET total_quantity = $1, average_price = $2, total_invested = $3, updated_at = now()
			WHERE user_asset_id = $4
			RETURNING updated_at
		`
		err = tx.QueryRow(ctx, holdingQuery, holding.TotalQuantity, holding.AveragePrice, holding.TotalInvested, holding.UserAssetID).
			Scan(&holding.UpdatedAt)
		if err != nil {
			slog.Error("failed to update holding", "error", err.Error())
			return util.NewInternalError("failed to update holding")
		}
	} else {
		holdingQuery := `
			INSERT INTO holdings (user_asset_id, total_quantity, average_price, total_invested)
			VALUES ($1, $2, $3, $4)
			RETURNING id, updated_at
		`
		err = tx.QueryRow(ctx, holdingQuery, holding.UserAssetID, holding.TotalQuantity, holding.AveragePrice, holding.TotalInvested).
			Scan(&holding.ID, &holding.UpdatedAt)
		if err != nil {
			slog.Error("failed to insert holding", "error", err.Error())
			return util.NewInternalError("failed to create holding")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		slog.Error("failed to commit transaction", "error", err.Error())
		return util.NewInternalError("failed to commit transaction")
	}

	return nil
}

func (r *TransactionRepository) GetAllByUserID(ctx context.Context, userID int64, limit, offset int) ([]dto.ResponseTransactionDto, error) {
	query := `
		SELECT t.id, t.user_asset_id, a.name, a.instrument_type, t.txn_type, t.quantity, t.price, t.txn_date
		FROM transactions t
		INNER JOIN user_assets ua ON t.user_asset_id = ua.id
		INNER JOIN assets a ON a.id = ua.asset_id
		WHERE ua.user_id = $1
		ORDER BY t.txn_date DESC, t.id DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		slog.Error("failed to list transactions", "error", err.Error())
		return nil, util.NewInternalError("failed to list transactions")
	}
	defer rows.Close()

	var transactions []dto.ResponseTransactionDto
	for rows.Next() {
		var txn dto.ResponseTransactionDto
		if err := rows.Scan(&txn.ID, &txn.UserAssetID, &txn.AssetName, &txn.AssetInstrumentType, &txn.TxnType, &txn.Quantity, &txn.Price, &txn.TxnDate); err != nil {
			slog.Error("failed to scan transaction row", "error", err.Error())
			return nil, util.NewInternalError("failed to list transactions")
		}
		transactions = append(transactions, txn)
	}

	if err := rows.Err(); err != nil {
		slog.Error("failed to iterate transaction rows", "error", err.Error())
		return nil, util.NewInternalError("failed to list transactions")
	}

	return transactions, nil
}

func (r *TransactionRepository) GetHoldingsByUserAssetID(ctx context.Context, userAssetID int64) (*model.Holding, error) {
	query := `
		SELECT id, user_asset_id, total_quantity, average_price, total_invested, updated_at
		FROM holdings
		WHERE user_asset_id = $1
	`

	var holding model.Holding
	err := r.db.QueryRow(ctx, query, userAssetID).
		Scan(&holding.ID, &holding.UserAssetID, &holding.TotalQuantity, &holding.AveragePrice, &holding.TotalInvested, &holding.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		slog.Error("failed to get holding by user asset id", "error", err.Error())
		return nil, util.NewInternalError("failed to get holding")
	}

	return &holding, nil
}

func (r *TransactionRepository) GetByID(ctx context.Context, id int64) (*model.Transaction, error) {
	query := `
		SELECT id, user_asset_id, txn_type, quantity, price, txn_date, created_at
		FROM transactions
		WHERE id = $1
	`

	var txn model.Transaction
	err := r.db.QueryRow(ctx, query, id).
		Scan(&txn.ID, &txn.UserAssetID, &txn.TxnType, &txn.Quantity, &txn.Price, &txn.TxnDate, &txn.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, util.NewNotFoundError(fmt.Sprintf("transaction with id %d not found", id))
		}
		slog.Error("failed to get transaction", "error", err.Error())
		return nil, util.NewInternalError("failed to get transaction")
	}

	return &txn, nil
}

func (r *TransactionRepository) Update(ctx context.Context, txn *model.Transaction) error {
	query := `
		UPDATE transactions
		SET txn_type = $1, quantity = $2, price = $3, txn_date = $4
		WHERE id = $5
	`

	res, err := r.db.Exec(ctx, query, txn.TxnType, txn.Quantity, txn.Price, txn.TxnDate, txn.ID)
	if err != nil {
		slog.Error("failed to update transaction", "error", err.Error())
		return util.NewInternalError("failed to update transaction")
	}

	if res.RowsAffected() == 0 {
		return util.NewNotFoundError(fmt.Sprintf("transaction with id %d not found", txn.ID))
	}

	return nil
}

func (r *TransactionRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM transactions WHERE id = $1`

	res, err := r.db.Exec(ctx, query, id)
	if err != nil {
		slog.Error("failed to delete transaction", "error", err.Error())
		return util.NewInternalError("failed to delete transaction")
	}

	if res.RowsAffected() == 0 {
		return util.NewNotFoundError(fmt.Sprintf("transaction with id %d not found", id))
	}

	return nil
}
