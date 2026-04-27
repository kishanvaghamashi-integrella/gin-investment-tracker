package service

import (
	"context"
	casparser "gin-investment-tracker/internal/cas-parser"
	repository "gin-investment-tracker/internal/repositories"
	"log/slog"
	"mime/multipart"
)

type StatementService struct {
	parser          casparser.CasParserInterface
	transactionRepo repository.TransactionRepositoryInterface
	holdingRepo     repository.HoldingRepositoryInterface
	userAssetRepo   repository.UserAssetRepositoryInterface
	statementRepo   repository.StatementRepositoryInterface
}

func NewCasStatementService(
	parser casparser.CasParserInterface,
	transactionRepo repository.TransactionRepositoryInterface,
	holdingRepo repository.HoldingRepositoryInterface,
	userAssetRepo repository.UserAssetRepositoryInterface,
	statementRepo repository.StatementRepositoryInterface,
) *StatementService {
	return &StatementService{
		parser:          parser,
		transactionRepo: transactionRepo,
		holdingRepo:     holdingRepo,
		userAssetRepo:   userAssetRepo,
		statementRepo:   statementRepo,
	}
}

func (s *StatementService) ProcessCasFile(ctx context.Context, file *multipart.FileHeader, filePassword string, userID int64) {
	bgCtx := context.WithoutCancel(ctx)
	go func() {
		casStatement, err := s.parser.ProcessCasFile(bgCtx, file, filePassword, userID)
		if err != nil {
			slog.Error("Failed to parse the file", "error", err.Error())
			return
		}
		slog.Info("CAS statement converted to JSON successfully.")

		if err := s.statementRepo.ProcessCASStatement(bgCtx, casStatement, userID); err != nil {
			slog.Error("Failed to process CAS statement", "error", err.Error())
			return
		}
		slog.Info("CAS statement uploaded successfully.", "userID", userID)
	}()
}
