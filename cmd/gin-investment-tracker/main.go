package main

import (
	"context"
	"gin-investment-tracker/internal/db"
	"gin-investment-tracker/internal/routes"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

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
	routes.RegisterRoutes(r, dbPool)
	r.Run(":8080")
}
