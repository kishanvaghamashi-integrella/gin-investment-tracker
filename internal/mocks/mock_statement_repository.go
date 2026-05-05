package mocks

import (
	"context"

	casparsermodel "gin-investment-tracker/internal/external-services/cas-parser/model"
	repository "gin-investment-tracker/internal/repositories"

	"github.com/stretchr/testify/mock"
)

type MockStatementRepository struct {
	mock.Mock
}

var _ repository.StatementRepositoryInterface = (*MockStatementRepository)(nil)

func (m *MockStatementRepository) ProcessCASStatement(ctx context.Context, cas *casparsermodel.CASStatement, userID int64) error {
	args := m.Called(ctx, cas, userID)
	return args.Error(0)
}
