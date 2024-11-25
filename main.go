package main

import (
	"log"
	"megajam/config"
	"megajam/gui"
	"megajam/logger"
)

func main() {
	// Initialize logger before loading configuration
	logFile, err := logger.InitLogger("logs/megajam.log")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logFile.Close()
	logger.Logger.Println("Starting -=* megaJAM *=-")

	// Load configuration early to ensure all components can access it if needed
	appConfig, err := config.LoadConfig("config/config.json")
	if err != nil {
		logger.Logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate configuration
	if err := config.ValidateConfig(appConfig); err != nil {
		logger.Logger.Fatalf("Configuration validation failed: %v", err)
	}

	// Start the GUI with the loaded configuration
	gui.CreateGUI(appConfig)
}
