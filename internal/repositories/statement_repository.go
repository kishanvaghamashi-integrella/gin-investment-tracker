package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	casparsermodel "gin-investment-tracker/internal/cas-parser/model"
	"gin-investment-tracker/internal/util"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StatementRepository struct {
	db *pgxpool.Pool
}

func NewStatementRepository(db *pgxpool.Pool) *StatementRepository {
	return &StatementRepository{db: db}
}

func (r *StatementRepository) ProcessCASStatement(ctx context.Context, cas *casparsermodel.CASStatement, userID int64) error {
	fromDate, err := parseDate(cas.StatementPeriod.From)
	if err != nil {
		return util.NewBadRequestError(fmt.Sprintf("invalid statement period from date: %s", cas.StatementPeriod.From))
	}
	toDate, err := parseDate(cas.StatementPeriod.To)
	if err != nil {
		return util.NewBadRequestError(fmt.Sprintf("invalid statement period to date: %s", cas.StatementPeriod.To))
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		slog.Error("failed to begin CAS import transaction", "error", err.Error())
		return util.NewInternalError("failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	for _, folio := range cas.Folios {
		for _, scheme := range folio.Schemes {
			if err := processScheme(ctx, tx, scheme, folio.AMC, userID, fromDate, toDate); err != nil {
				return err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		slog.Error("failed to commit CAS import transaction", "error", err.Error())
		return util.NewInternalError("failed to commit transaction")
	}

	return nil
}

func processScheme(ctx context.Context, tx pgx.Tx, scheme casparsermodel.Scheme, amc string, userID int64, fromDate, toDate time.Time) error {
	assetID, err := findOrCreateAsset(ctx, tx, scheme, amc)
	if err != nil {
		return err
	}

	userAssetID, err := findOrCreateUserAsset(ctx, tx, userID, assetID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`DELETE FROM transactions WHERE user_asset_id = $1 AND txn_date BETWEEN $2 AND $3`,
		userAssetID, fromDate, toDate,
	)
	if err != nil {
		slog.Error("failed to delete transactions in date range", "error", err.Error())
		return util.NewInternalError("failed to delete transactions")
	}

	for _, txnData := range scheme.Transactions {
		txnType := mapTxnType(txnData.Type)
		if txnType == "" {
			continue
		}

		txnDate, err := parseDate(txnData.Date)
		if err != nil {
			slog.Warn("skipping transaction with unparseable date", "date", txnData.Date)
			continue
		}

		quantity, err := strconv.ParseFloat(txnData.Units, 64)
		if err != nil || quantity == 0 {
			slog.Warn("skipping transaction with invalid units", "units", txnData.Units)
			continue
		}

		price, err := strconv.ParseFloat(txnData.NAV, 64)
		if err != nil {
			slog.Warn("skipping transaction with invalid nav", "nav", txnData.NAV)
			continue
		}

		var desc *string
		if txnData.Description != "" {
			desc = &txnData.Description
		}

		_, err = tx.Exec(ctx,
			`INSERT INTO transactions (user_asset_id, txn_type, quantity, price, txn_date, description)
			 VALUES ($1, $2, $3, $4, $5, $6)`,
			userAssetID, txnType, quantity, price, txnDate, desc,
		)
		if err != nil {
			slog.Error("failed to insert transaction from CAS data", "error", err.Error())
			return util.NewInternalError("failed to insert transaction")
		}
	}

	return recalculateHolding(ctx, tx, userAssetID)
}

func findOrCreateAsset(ctx context.Context, tx pgx.Tx, scheme casparsermodel.Scheme, amc string) (int64, error) {
	var assetID int64
	err := tx.QueryRow(ctx, `SELECT id FROM assets WHERE isin = $1`, scheme.ISIN).Scan(&assetID)
	if err == nil {
		return assetID, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("failed to query asset by isin", "error", err.Error())
		return 0, util.NewInternalError("failed to query asset")
	}

	var amcPtr *string
	if amc != "" {
		amcPtr = &amc
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO assets (symbol, name, instrument_type, isin, amc, exchange, currency, external_platform_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id`,
		scheme.ISIN, scheme.Scheme, "mutual_fund", scheme.ISIN, amcPtr, "MF", "INR", scheme.AMFI,
	).Scan(&assetID)
	if err != nil {
		slog.Error("failed to create asset from CAS data", "error", err.Error())
		return 0, util.NewInternalError("failed to create asset")
	}

	return assetID, nil
}

func findOrCreateUserAsset(ctx context.Context, tx pgx.Tx, userID, assetID int64) (int64, error) {
	var userAssetID int64
	err := tx.QueryRow(ctx,
		`SELECT id FROM user_assets WHERE user_id = $1 AND asset_id = $2`,
		userID, assetID,
	).Scan(&userAssetID)
	if err == nil {
		return userAssetID, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("failed to query user_asset", "error", err.Error())
		return 0, util.NewInternalError("failed to query user asset")
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO user_assets (user_id, asset_id) VALUES ($1, $2) RETURNING id`,
		userID, assetID,
	).Scan(&userAssetID)
	if err != nil {
		slog.Error("failed to create user_asset", "error", err.Error())
		return 0, util.NewInternalError("failed to create user asset")
	}

	return userAssetID, nil
}

func recalculateHolding(ctx context.Context, tx pgx.Tx, userAssetID int64) error {
	var totalQty, avgPrice, totalInvested float64
	err := tx.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(CASE WHEN txn_type = 'BUY' THEN quantity ELSE -quantity END), 0),
			CASE
				WHEN SUM(CASE WHEN txn_type = 'BUY' THEN quantity ELSE 0 END) > 0
				THEN SUM(CASE WHEN txn_type = 'BUY' THEN quantity * price ELSE 0 END) /
				     SUM(CASE WHEN txn_type = 'BUY' THEN quantity ELSE 0 END)
				ELSE 0
			END,
			COALESCE(
				SUM(CASE WHEN txn_type = 'BUY'  THEN quantity * price ELSE 0 END) -
				SUM(CASE WHEN txn_type = 'SELL' THEN quantity * price ELSE 0 END),
			0)
		FROM transactions
		WHERE user_asset_id = $1
	`, userAssetID).Scan(&totalQty, &avgPrice, &totalInvested)
	if err != nil {
		slog.Error("failed to compute holding aggregates", "error", err.Error())
		return util.NewInternalError("failed to recalculate holdings")
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO holdings (user_asset_id, total_quantity, average_price, total_invested)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_asset_id) DO UPDATE
			SET total_quantity = $2,
			    average_price  = $3,
			    total_invested = $4,
			    updated_at     = now()
	`, userAssetID, totalQty, avgPrice, totalInvested)
	if err != nil {
		slog.Error("failed to upsert holding", "error", err.Error())
		return util.NewInternalError("failed to upsert holding")
	}

	return nil
}

func mapTxnType(casType string) string {
	switch casType {
	case "PURCHASE", "SWITCH_IN", "PURCHASE_SIP":
		return "BUY"
	case "REDEMPTION", "SWITCH_OUT":
		return "SELL"
	default:
		return ""
	}
}

func parseDate(s string) (time.Time, error) {
	formats := []string{
		"2006-01-02",
		"02-Jan-2006",
		"02-01-2006",
		"Jan 02, 2006",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unparseable date: %s", s)
}
