package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/transactions-platform/docs" // Import swagger docs
	"github.com/transactions-platform/internal/database"
	"github.com/transactions-platform/internal/handlers"
	"github.com/transactions-platform/internal/logger"
	"github.com/transactions-platform/internal/repository"
	"github.com/transactions-platform/internal/service"
)

type API struct {
	server *http.Server
	router *gin.Engine
	db     *sql.DB
	port   string
}

func Build(ctx context.Context) (*API, error) {
	// Initialize logger
	logger.Init()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Starting application").
		Str("port", port).
		Str("environment", os.Getenv("APP_ENV")).
		Send()

	// Connect to database
	dbConfig := database.NewConfigFromEnv()
	db, err := database.Connect(dbConfig)
	if err != nil {
		logger.Error("Failed to connect to database").Err(err).Send()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	logger.Info("Database connection established").Send()

	// Set Gin mode based on environment
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(logger.GinLogger())
	router.Use(logger.GinRecovery())

	// Initialize repositories
	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// Initialize services
	accountService := service.NewAccountService(accountRepo)
	transactionService := service.NewTransactionService(transactionRepo, accountRepo)

	// Initialize handlers
	accountHandler := handlers.NewAccountHandler(accountService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Register health endpoint
	router.GET("/health", handlers.HealthCheck)

	// Register account endpoints
	router.POST("/accounts", accountHandler.CreateAccount)
	router.GET("/accounts/:id", accountHandler.GetAccount)

	// Register transaction endpoints
	router.POST("/transactions", transactionHandler.CreateTransaction)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &API{
		server: server,
		router: router,
		db:     db,
		port:   port,
	}, nil
}

func (a *API) Run(ctx context.Context) error {
	errChan := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		logger.Info("Server listening").
			Str("port", a.port).
			Str("address", a.server.Addr).
			Send()
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start").Err(err).Send()
			errChan <- fmt.Errorf("server failed to start: %w", err)
		}
	}()

	// Wait for context cancellation or server error
	select {
	case <-ctx.Done():
		logger.Info("Shutdown signal received, gracefully shutting down server").Send()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := a.server.Shutdown(shutdownCtx); err != nil {
			logger.Error("Server forced to shutdown").Err(err).Send()
			return fmt.Errorf("server forced to shutdown: %w", err)
		}
		logger.Info("Server shutdown completed successfully").Send()
		return nil
	case err := <-errChan:
		return err
	}
}

func (a *API) Close(ctx context.Context) error {
	logger.Info("Closing application resources").Send()

	// Close database connection
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			logger.Error("Error closing database connection").Err(err).Send()
		} else {
			logger.Info("Database connection closed").Send()
		}
	}

	// Close server
	if a.server != nil {
		return a.server.Close()
	}
	return nil
}
