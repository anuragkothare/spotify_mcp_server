package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anuragkothare/spotify_mcp_server/internal/config"
	"github.com/anuragkothare/spotify_mcp_server/internal/handlers"
	"github.com/anuragkothare/spotify_mcp_server/internal/mcp"
	"github.com/anuragkothare/spotify_mcp_server/internal/spotify"
	"github.com/anuragkothare/spotify_mcp_server/pkg/logger"
)

func main() {
	// Initialize logger
	log := logger.New()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Spotify client
	spotifyClient, err := spotify.NewClient(cfg.Spotify)
	if err != nil {
		log.Fatalf("Failed to create Spotify client: %v", err)
	}

	// Initialize MCP server
	mcpServer := mcp.NewServer(spotifyClient, log)

	// Initialize HTTP handlers
	handler := handlers.NewHandler(mcpServer, log)

	// Setup HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/mcp", handler.HandleMCP)
	mux.HandleFunc("/health", handler.HandleHealth)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      mux,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Infof("Starting MCP server on port %d", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Info("Server exited")
}
