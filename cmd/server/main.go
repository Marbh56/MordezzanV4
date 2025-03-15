package main

import (
	"context"
	"fmt"
	"mordezzanV4/internal/app"
	"mordezzanV4/internal/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		// Not a fatal error, as .env might not exist in production
		fmt.Println("Warning: .env file not found")
	}
	// Initialize the Zap logger
	logger.Init(logger.Config{
		LogLevel:         logger.LogLevelDebug,
		IncludeTimestamp: true,
		IncludeFileLine:  true,
		Development:      true, // Set to false in production
	})

	logger.Info("Starting server...")

	// Create a new application instance
	app, err := app.NewApp("./myproject.db")
	if err != nil {
		logger.Fatal("Failed to initialize application: %v", err)
	}
	defer app.Shutdown()

	// Set up the HTTP routes
	handler := app.SetupRoutes()

	// Get the port from environment variables or use the default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Configure the HTTP server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start the server in a goroutine
	go func() {
		logger.Info("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start: %v", err)
		}
	}()

	// Set up graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logger.Info("Shutting down server...")

	// Create a deadline for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: %v", err)
	}

	logger.Info("Server stopped successfully")
}
