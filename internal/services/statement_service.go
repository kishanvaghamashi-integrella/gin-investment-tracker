package service

import (
	"context"
	casparser "gin-investment-tracker/internal/external-services/cas-parser"
	repository "gin-investment-tracker/internal/repositories"
	"gin-investment-tracker/internal/util"
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
		log := util.FromContext(bgCtx).With("service", "StatementService.ProcessCasFile", "user_id", userID)

		casStatement, err := s.parser.ProcessCasFile(bgCtx, file, filePassword, userID)
		if err != nil {
			log.Errorw("Failed to parse the file", "error", err)
			return
		}
		log.Infow("CAS statement converted to JSON successfully.")

		if err := s.statementRepo.ProcessCASStatement(bgCtx, casStatement, userID); err != nil {
			log.Errorw("Failed to process CAS statement", "error", err)
			return
		}
		log.Infow("CAS statement uploaded successfully.")
	}()
}
