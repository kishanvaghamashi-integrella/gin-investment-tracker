package main

import (
	"context"
	_ "gin-investment-tracker/docs"
	"gin-investment-tracker/internal/db"
	"gin-investment-tracker/internal/server"
	"log"
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
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	if env := godotenv.Load(); env != nil {
		log.Fatal("Error loading .env file")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbPool, err := db.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error while connecting DB %s", err.Error())
		return
	}

	r := gin.Default()
	server.RegisterRoutes(r, dbPool)
	r.Run(":8080")
}
