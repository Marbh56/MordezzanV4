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
		fmt.Println("Warning: .env file not found")
	}

	logger.Init(logger.Config{
		LogLevel:         logger.LogLevelDebug,
		IncludeTimestamp: true,
		IncludeFileLine:  true,
		Development:      true,
	})

	logger.Info("Starting server...")
	app, err := app.NewApp("./myproject.db")
	if err != nil {
		logger.Fatal("Failed to initialize application: %v", err)
	}
	defer app.Shutdown()

	handler := app.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Use localhost instead of local IP address
	addr := "localhost:" + port

	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("Server starting on http://%s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: %v", err)
	}

	logger.Info("Server stopped successfully")
}
