package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"trading-platform/backend/internal/portfolioanalytics"
	"trading-platform/backend/internal/orderexecution"
	"trading-platform/backend/internal/auth"
	"trading-platform/backend/internal/api"
	"trading-platform/backend/internal/websocket"
	"trading-platform/backend/internal/marketdata"

	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "TRADING-PLATFORM: ", log.LstdFlags|log.Lshortfile)
	logger.Println("Starting trading platform...")

	// Load configuration
	// In a real implementation, this would load from a config file or environment variables
	dbConnStr := "postgres://postgres:postgres@localhost:5432/tradingplatform?sslmode=disable"
	serverAddr := ":8080"
	
	// Connect to database
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		logger.Fatalf("Failed to ping database: %v", err)
	}
	logger.Println("Connected to database")

	// Initialize repositories
	portfolioRepo := portfolioanalytics.NewPostgresRepository(db)
	
	// Initialize market data provider
	marketDataProvider := marketdata.NewMarketDataProvider()
	
	// Initialize portfolio analytics engine
	analyticsEngine := portfolioanalytics.NewPortfolioAnalyticsEngine(marketDataProvider, 5)
	if err := analyticsEngine.Start(); err != nil {
		logger.Fatalf("Failed to start analytics engine: %v", err)
	}
	defer analyticsEngine.Stop()
	
	// Initialize services
	portfolioService := portfolioanalytics.NewService(portfolioRepo, analyticsEngine)
	orderExecutionService := orderexecution.NewService(db)
	authService := auth.NewService(db)
	
	// Initialize controllers
	portfolioController := portfolioanalytics.NewController(portfolioService)
	orderExecutionController := orderexecution.NewController(orderExecutionService)
	authController := auth.NewController(authService)
	
	// Initialize WebSocket handler
	wsHandler := websocket.NewHandler(portfolioService, orderExecutionService)
	
	// Initialize router
	router := mux.NewRouter()
	
	// Register API routes
	api.RegisterRoutes(router, portfolioController, orderExecutionController, authController)
	
	// Register WebSocket handler
	router.HandleFunc("/ws", wsHandler.HandleConnection)
	
	// Create HTTP server
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	// Start server in a goroutine
	go func() {
		logger.Printf("Server listening on %s", serverAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()
	
	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Println("Shutting down server...")
	
	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}
	
	logger.Println("Server exited properly")
}
