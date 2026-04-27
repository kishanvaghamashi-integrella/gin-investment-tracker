package repository

import (
	"context"

	casparsermodel "gin-investment-tracker/internal/cas-parser/model"
)

type StatementRepositoryInterface interface {
	ProcessCASStatement(ctx context.Context, cas *casparsermodel.CASStatement, userID int64) error
}
