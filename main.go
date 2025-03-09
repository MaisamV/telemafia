package main

import (
	"log"
	"telemafia/config"
)

func main() {
	// Load Configuration
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Dependencies
	botHandler, err := config.InitializeDependencies(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Register bot handlers
	botHandler.RegisterHandlers()

	log.Println("Bot is running...")
	botHandler.Start()
}
