package service

import (
	"context"
	"mime/multipart"
)

type StatementServiceInterface interface {
	ProcessCasFile(ctx context.Context, file *multipart.FileHeader, filePassword string, userID int64)
}
