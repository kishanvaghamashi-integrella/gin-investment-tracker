package main

import (
	"context"
	_ "gin-investment-tracker/docs"
	"gin-investment-tracker/internal/db"
	"gin-investment-tracker/internal/server"
	"gin-investment-tracker/internal/util"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title MF Stock Tracker API
// @version 1.0
// @description API for managing users and assets in MF Stock Tracker.
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name jwt_token
// @description JWT token stored in an HttpOnly cookie. Set automatically on login.

func main() {
	if env := godotenv.Load(); env != nil {
		util.Logger.Errorw("failed to load .env file", "error", env)
		os.Exit(1)
		return
	}

	util.InitLogger()
	defer util.SyncLogger()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbPool, err := db.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		util.Logger.Errorw("failed while connecting db", "error", err)
		os.Exit(1)
		return
	}

	r := gin.Default()
	if err := server.RegisterRoutes(r, dbPool); err != nil {
		util.Logger.Errorw("failed to register routes", "error", err)
		os.Exit(1)
		return
	}

	if err := r.Run(":8000"); err != nil {
		util.Logger.Errorw("server failed to run", "error", err)
		os.Exit(1)
	}
}
