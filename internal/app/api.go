package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/transactions-platform/docs" // Import swagger docs
	"github.com/transactions-platform/internal/database"
	"github.com/transactions-platform/internal/handlers"
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Connect to database
	dbConfig := database.NewConfigFromEnv()
	db, err := database.Connect(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set Gin mode based on environment
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

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
		log.Printf("Starting server on port %s", a.port)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("server failed to start: %w", err)
		}
	}()

	// Wait for context cancellation or server error
	select {
	case <-ctx.Done():
		log.Println("Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := a.server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server forced to shutdown: %w", err)
		}
		log.Println("Server exited gracefully")
		return nil
	case err := <-errChan:
		return err
	}
}

func (a *API) Close(ctx context.Context) error {
	// Close database connection
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}

	// Close server
	if a.server != nil {
		return a.server.Close()
	}
	return nil
}
