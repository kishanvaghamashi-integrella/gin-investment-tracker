package repository

import (
	"context"
	casparsermodel "gin-investment-tracker/internal/external-services/cas-parser/model"
)

type StatementRepositoryInterface interface {
	ProcessCASStatement(ctx context.Context, cas *casparsermodel.CASStatement, userID int64) error
}
