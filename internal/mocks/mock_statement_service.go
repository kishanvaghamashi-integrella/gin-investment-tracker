package mocks

import (
	"context"
	"mime/multipart"

	service "gin-investment-tracker/internal/services"

	"github.com/stretchr/testify/mock"
)

type MockStatementService struct {
	mock.Mock
}

var _ service.StatementServiceInterface = (*MockStatementService)(nil)

func (m *MockStatementService) ProcessCasFile(ctx context.Context, file *multipart.FileHeader, filePassword string, userID int64) {
	m.Called(ctx, file, filePassword, userID)
}
