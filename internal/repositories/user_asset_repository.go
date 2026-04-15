package repository

import (
	"context"
	"fmt"
	model "gin-investment-tracker/internal/models"
	"gin-investment-tracker/internal/util"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserAssetRepository struct {
	db *pgxpool.Pool
}

func NewUserAssetRepository(db *pgxpool.Pool) *UserAssetRepository {
	return &UserAssetRepository{db: db}
}

func (r *UserAssetRepository) Create(ctx context.Context, userAsset *model.UserAsset) error {
	query := `
		INSERT INTO user_assets (user_id, asset_id)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query, userAsset.UserID, userAsset.AssetID).
		Scan(&userAsset.ID, &userAsset.CreatedAt)
	if err != nil {
		slog.Error("failed to create user asset", "error", err.Error())
		return util.NewInternalError("failed to create user asset")
	}

	return nil
}

func (r *UserAssetRepository) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]model.UserAsset, error) {
	query := `
		SELECT id, user_id, asset_id, created_at
		FROM user_assets
		WHERE user_id = $1
		ORDER BY id
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		slog.Error("failed to list user assets", "error", err.Error())
		return nil, util.NewInternalError("failed to list user assets")
	}
	defer rows.Close()

	var userAssets []model.UserAsset
	for rows.Next() {
		var ua model.UserAsset
		if err := rows.Scan(&ua.ID, &ua.UserID, &ua.AssetID, &ua.CreatedAt); err != nil {
			slog.Error("failed to scan user asset row", "error", err.Error())
			return nil, util.NewInternalError("failed to list user assets")
		}
		userAssets = append(userAssets, ua)
	}

	if err := rows.Err(); err != nil {
		slog.Error("failed during user asset row iteration", "error", err.Error())
		return nil, util.NewInternalError("failed to list user assets")
	}

	return userAssets, nil
}

func (r *UserAssetRepository) Delete(ctx context.Context, id, userID int64) error {
	query := `DELETE FROM user_assets WHERE id = $1 AND user_id = $2`

	res, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		slog.Error("failed to delete user asset", "error", err.Error())
		return util.NewInternalError("failed to delete user asset")
	}

	if res.RowsAffected() == 0 {
		return util.NewNotFoundError(fmt.Sprintf("user asset with id %d not found for user %d", id, userID))
	}

	return nil
}

func (r *UserAssetRepository) IsUserAssetExists(ctx context.Context, userID int64, assetID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_assets WHERE user_id = $1 AND asset_id = $2)`

	var exists bool
	if err := r.db.QueryRow(ctx, query, userID, assetID).Scan(&exists); err != nil {
		slog.Error("failed to check asset existence", "error", err.Error())
		return false, util.NewInternalError("failed to check asset existence")
	}

	return exists, nil
}

func (r *UserAssetRepository) ExistsByID(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_assets WHERE id = $1)`

	var exists bool
	if err := r.db.QueryRow(ctx, query, id).Scan(&exists); err != nil {
		slog.Error("failed to check user asset existence", "error", err.Error())
		return false, util.NewInternalError("failed to check user asset existence")
	}

	return exists, nil
}

func (r *UserAssetRepository) GetIdByUserIdAssetId(ctx context.Context, userID, assetID int64) (*int64, error) {
	query := `SELECT id FROM user_assets WHERE user_id=$1 AND asset_id = $2`

	var userAssetId int64
	if err := r.db.QueryRow(ctx, query, userID, assetID).Scan(&userAssetId); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		slog.Error("failed to retrieve userAssetID", "error", err.Error())
		return nil, util.NewInternalError("failed to retrieve userAssetID")
	}
	return &userAssetId, nil
}
