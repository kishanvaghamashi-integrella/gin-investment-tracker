package routes

import (
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
	// Repositories
	userRepository := repository.NewUserRepository(db)
	assetRepository := repository.NewAssetRepository(db)
	userAssetRepository := repository.NewUserAssetRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)

	// Services
	userService := service.NewUserService(userRepository)
	assetService := service.NewAssetService(assetRepository)
	userAssetService := service.NewUserAssetService(userAssetRepository, userRepository, assetRepository)
	transactionService := service.NewTransactionService(transactionRepository, userAssetRepository, userRepository, assetRepository)

	// Handlers
	userHandler := handler.NewUserHandler(userService)
	assetHandler := handler.NewAssetHandler(assetService)
	userAssetHandler := handler.NewUserAssetHandler(userAssetService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// routes
	if isDevelopmentEnvironment() {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	unprotectedRouter := r.Group("/api")
	userHandler.SetRoutes(unprotectedRouter)

	protectedRouter := r.Group("/api")
	protectedRouter.Use(middleware.JWTAuth())
	{
		assetHandler.SetRoutes(protectedRouter)
		userAssetHandler.SetRoutes(protectedRouter)
		transactionHandler.SetRoutes(protectedRouter)
	}
}

func isDevelopmentEnvironment() bool {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	return env == "dev" || env == "development" || env == "local"
}
