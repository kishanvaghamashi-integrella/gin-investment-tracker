package routes

import (
	handler "gin-investment-tracker/internal/handlers"
	repositoryimpl "gin-investment-tracker/internal/repositories_impl"
	service "gin-investment-tracker/internal/services"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "gin-investment-tracker/docs"
)

func RegisterRoutes(r *gin.Engine, db *pgxpool.Pool) {
	// Repositories
	userRepository := repositoryimpl.NewUserRepository(db)

	// Services
	userService := service.NewUserService(userRepository)

	// Handlers
	userHandler := handler.NewUserHandler(userService)

	// routes
	if isDevelopmentEnvironment() {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	unprotectedRouter := r.Group("/api")
	userHandler.SetRoutes(unprotectedRouter)
}

func isDevelopmentEnvironment() bool {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	return env == "dev" || env == "development" || env == "local"
}
