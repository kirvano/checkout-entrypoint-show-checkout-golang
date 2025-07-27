package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"checkout-go/internal/handlers"
	"checkout-go/internal/infrastructure/di"
)

func main() {
	log.Println("Starting Gin-based checkout server...")

	// Initialize dependency injection container
	log.Println("Initializing dependency injection container...")
	container, err := di.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize DI container: %v", err)
	}
	log.Println("DI container initialized successfully")

	// Get configuration from container
	config := container.GetConfig()

	// Test that we can get the use case
	useCase := container.GetShowCheckoutUseCase()
	if useCase == nil {
		log.Fatalf("Failed to get ShowCheckoutUseCase from container")
	}
	log.Println("ShowCheckoutUseCase retrieved successfully")

	// Set gin mode based on configuration
	gin.SetMode(config.GetGinMode())

	// Create Gin router
	router := gin.New()

	// Add middleware
	router.Use(handlers.LoggerMiddleware())
	router.Use(handlers.RecoveryMiddleware())
	router.Use(handlers.CORSMiddleware())
	router.Use(handlers.RequestIDMiddleware())

	// Initialize handlers
	checkoutHandlers := handlers.NewCheckoutHandlers(container)

	// Define routes
	setupRoutes(router, checkoutHandlers)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", config.Port)
		log.Printf("Environment: %s", config.AppEnv)
		log.Printf("Gin mode: %s", config.GetGinMode())
		
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

// setupRoutes configures the API routes
func setupRoutes(router *gin.Engine, handlers *handlers.CheckoutHandlers) {
	// Health check endpoint
	router.GET("/health", handlers.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Checkout routes
		checkout := v1.Group("/checkout")
		{
			checkout.GET("/:uuid", handlers.ShowCheckout)
		}
	}

	// Root checkout route for backward compatibility
	router.GET("/checkout/:uuid", handlers.ShowCheckout)
} 