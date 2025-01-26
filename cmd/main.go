package main

import (
	"log"

	"github.com/NBDor/eternalsphere-auth/internal/config"
	"github.com/NBDor/eternalsphere-auth/internal/handlers"
	"github.com/NBDor/eternalsphere-auth/internal/repository/postgres"
	"github.com/NBDor/eternalsphere-auth/internal/service"
	shared "github.com/NBDor/eternalsphere-shared-go/database/postgres"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// Initialize database
	dbConfig := shared.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
	}

	conn, err := shared.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close()

	// Initialize repositories
	userRepo := postgres.NewUserRepository(conn)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.JWT.SecretKey)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Setup router
	router := gin.Default()

	// Routes
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)
	router.POST("/refresh", authHandler.RefreshToken)

	// Start server
	if err := router.Run(cfg.Server.Address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
