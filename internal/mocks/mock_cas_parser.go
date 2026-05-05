package mocks

import (
	"context"

	casparser "gin-investment-tracker/internal/external-services/cas-parser"
	casparsermodel "gin-investment-tracker/internal/external-services/cas-parser/model"
	"mime/multipart"

	"github.com/stretchr/testify/mock"
)

type MockCasParser struct {
	mock.Mock
}

var _ casparser.CasParserInterface = (*MockCasParser)(nil)

func (m *MockCasParser) ProcessCasFile(ctx context.Context, file *multipart.FileHeader, filePassword string, userID int64) (*casparsermodel.CASStatement, error) {
	args := m.Called(ctx, file, filePassword, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*casparsermodel.CASStatement), args.Error(1)
}
