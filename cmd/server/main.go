package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ReesavGupta/discord-mcp-server/internal/config"
	"github.com/ReesavGupta/discord-mcp-server/internal/mcp"
	"github.com/sirupsen/logrus"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup logger
	logger := logrus.New()
	logger.SetOutput(os.Stderr)

	level, err := logrus.ParseLevel(cfg.Logging.Level)
	if err != nil {
		logger.Warn("Invalid log level, using info")
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	if cfg.Logging.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	// Create and start server
	server, err := mcp.NewServer(cfg, logger)
	if err != nil {
		logger.Fatalf("Failed to create server: %v", err)
	}

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		logger.Info("Shutting down server...")
		os.Exit(0)
	}()

	// Start server
	if err := server.Start(); err != nil {
		logger.Fatalf("Server failed: %v", err)
	}
}
