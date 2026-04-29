package casparser

import (
	"context"
	casparsermodel "gin-investment-tracker/internal/external-services/cas-parser/model"
	"mime/multipart"
)

type CasParserInterface interface {
	ProcessCasFile(ctx context.Context, file *multipart.FileHeader, filePassword string, userID int64) (*casparsermodel.CASStatement, error)
}
