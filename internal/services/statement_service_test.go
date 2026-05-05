package service_test

import (
	"context"
	"errors"
	"mime/multipart"
	"sync"
	"testing"
	"time"

	casparsermodel "gin-investment-tracker/internal/external-services/cas-parser/model"
	"gin-investment-tracker/internal/mocks"
	service "gin-investment-tracker/internal/services"

	"github.com/stretchr/testify/mock"
)

// waitForGoroutine waits up to 3 seconds for a WaitGroup to reach zero.
func waitForGoroutine(t *testing.T, wg *sync.WaitGroup) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("background goroutine did not complete within 3 seconds")
	}
}

func newStatementService(
	parser *mocks.MockCasParser,
	txnRepo *mocks.MockTransactionRepository,
	holdingRepo *mocks.MockHoldingRepository,
	userAssetRepo *mocks.MockUserAssetRepository,
	statementRepo *mocks.MockStatementRepository,
) *service.StatementService {
	return service.NewCasStatementService(parser, txnRepo, holdingRepo, userAssetRepo, statementRepo)
}

// ─────────────────────────────────────────────
// ProcessCasFile
// ─────────────────────────────────────────────

func TestStatementService_ProcessCasFile_Success(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	parser := new(mocks.MockCasParser)
	txnRepo := new(mocks.MockTransactionRepository)
	holdingRepo := new(mocks.MockHoldingRepository)
	userAssetRepo := new(mocks.MockUserAssetRepository)
	statementRepo := new(mocks.MockStatementRepository)

	svc := newStatementService(parser, txnRepo, holdingRepo, userAssetRepo, statementRepo)

	casStatement := &casparsermodel.CASStatement{}
	fileHeader := &multipart.FileHeader{Filename: "statement.pdf"}

	parser.On("ProcessCasFile", mock.Anything, fileHeader, "pass", int64(1)).
		Return(casStatement, nil).
		Run(func(_ mock.Arguments) { wg.Done() })
	statementRepo.On("ProcessCASStatement", mock.Anything, casStatement, int64(1)).Return(nil)

	svc.ProcessCasFile(context.Background(), fileHeader, "pass", 1)

	waitForGoroutine(t, &wg)
	parser.AssertExpectations(t)
	statementRepo.AssertExpectations(t)
}

func TestStatementService_ProcessCasFile_ParserError_DoesNotCallRepo(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	parser := new(mocks.MockCasParser)
	txnRepo := new(mocks.MockTransactionRepository)
	holdingRepo := new(mocks.MockHoldingRepository)
	userAssetRepo := new(mocks.MockUserAssetRepository)
	statementRepo := new(mocks.MockStatementRepository)

	svc := newStatementService(parser, txnRepo, holdingRepo, userAssetRepo, statementRepo)

	fileHeader := &multipart.FileHeader{Filename: "bad.pdf"}

	parser.On("ProcessCasFile", mock.Anything, fileHeader, "", int64(1)).
		Return(nil, errors.New("parse error")).
		Run(func(_ mock.Arguments) { wg.Done() })

	svc.ProcessCasFile(context.Background(), fileHeader, "", 1)

	waitForGoroutine(t, &wg)
	parser.AssertExpectations(t)
	statementRepo.AssertNotCalled(t, "ProcessCASStatement")
}

func TestStatementService_ProcessCasFile_RepoError(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	parser := new(mocks.MockCasParser)
	txnRepo := new(mocks.MockTransactionRepository)
	holdingRepo := new(mocks.MockHoldingRepository)
	userAssetRepo := new(mocks.MockUserAssetRepository)
	statementRepo := new(mocks.MockStatementRepository)

	svc := newStatementService(parser, txnRepo, holdingRepo, userAssetRepo, statementRepo)

	casStatement := &casparsermodel.CASStatement{}
	fileHeader := &multipart.FileHeader{Filename: "statement.pdf"}

	parser.On("ProcessCasFile", mock.Anything, fileHeader, "", int64(1)).
		Return(casStatement, nil).
		Run(func(_ mock.Arguments) { wg.Done() })
	statementRepo.On("ProcessCASStatement", mock.Anything, casStatement, int64(1)).
		Return(errors.New("db failure"))

	svc.ProcessCasFile(context.Background(), fileHeader, "", 1)

	waitForGoroutine(t, &wg)
	parser.AssertExpectations(t)
	statementRepo.AssertExpectations(t)
}
