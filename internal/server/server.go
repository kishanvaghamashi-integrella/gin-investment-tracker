package server

import (
	"gin-investment-tracker/internal/cron"
	assetprice "gin-investment-tracker/internal/external-services/asset-details"
	casparser "gin-investment-tracker/internal/external-services/cas-parser"
	handler "gin-investment-tracker/internal/handlers"
	middleware "gin-investment-tracker/internal/middlewares"
	repository "gin-investment-tracker/internal/repositories"
	service "gin-investment-tracker/internal/services"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine, db *pgxpool.Pool) {
	// 3rd party services
	casParser := casparser.NewCasParserPythonApi()
	mfPriceFetcher := assetprice.NewMfapiFetcher()
	stockPriceFetcher := assetprice.NewYahooStockFetcher()
	assetPriceFetcher := assetprice.NewAssetPriceService(stockPriceFetcher, mfPriceFetcher)

	// Repositories
	userRepository := repository.NewUserRepository(db)
	assetRepository := repository.NewAssetRepository(db)
	userAssetRepository := repository.NewUserAssetRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)
	holdingRepository := repository.NewHoldingRepository(db)
	statementRepository := repository.NewStatementRepository(db)
	priceDetailRepository := repository.NewPriceDetailRepository(db)
	dashboardRepository := repository.NewDashboardRepository(db)

	// Services
	userService := service.NewUserService(userRepository)
	assetService := service.NewAssetService(assetRepository)
	userAssetService := service.NewUserAssetService(userAssetRepository, userRepository, assetRepository)
	transactionService := service.NewTransactionService(transactionRepository, userAssetRepository, userRepository, assetRepository)
	holdingService := service.NewHoldingService(holdingRepository, userRepository)
	casStatementService := service.NewCasStatementService(casParser, transactionRepository, holdingRepository, userAssetRepository, statementRepository)
	dashboardService := service.NewDashboardService(dashboardRepository)

	// Handlers
	userHandler := handler.NewUserHandler(userService)
	assetHandler := handler.NewAssetHandler(assetService)
	userAssetHandler := handler.NewUserAssetHandler(userAssetService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	holdingHandler := handler.NewHoldingHandler(holdingService)
	casStatementHandler := handler.NewCasStatementHandler(casStatementService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)

	// Cron Job
	cronJob := cron.NewCronJobs(assetRepository, priceDetailRepository, assetPriceFetcher)
	cronJob.Start()

	// routes
	if isDevelopmentEnvironment() {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	r.Use(CORSMiddleware())

	unprotectedRouter := r.Group("/api")
	userHandler.SetRoutes(unprotectedRouter)

	protectedRouter := r.Group("/api")
	protectedRouter.Use(middleware.JWTAuth())
	{
		assetHandler.SetRoutes(protectedRouter)
		userAssetHandler.SetRoutes(protectedRouter)
		transactionHandler.SetRoutes(protectedRouter)
		holdingHandler.SetRoutes(protectedRouter)
		casStatementHandler.SetRoutes(protectedRouter)
		dashboardHandler.SetRoutes(protectedRouter)
	}
}

func isDevelopmentEnvironment() bool {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	return env == "dev" || env == "development" || env == "local"
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
