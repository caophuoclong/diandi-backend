/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	"log"
	"os"

	"diandi-backend/api/handlers"
	"diandi-backend/config"
	"diandi-backend/repositories"
	"diandi-backend/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Connect to MongoDB
	ctx := context.Background()
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database(os.Getenv("MONGODB_DATABASE"))

	// Initialize repositories
	oauthRepo := repositories.NewMongoOAuthRepository(db)

	// Load OAuth configurations
	oauthConfigs := config.LoadOAuthConfigs()

	// Initialize services
	oauthService := services.NewOAuthService(oauthRepo, oauthConfigs.GetAllConfigs())

	// Initialize handlers
	oauthHandler := handlers.NewOAuthHandler(oauthService)

	// Setup Gin router
	router := gin.Default()

	// API routes
	v1 := router.Group("/api/v1")
	{
		oauthHandler.RegisterRoutes(v1)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
