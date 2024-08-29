package main

import (
	"log"

	"pairs-trading-backend/internal/config"
	"pairs-trading-backend/internal/database"
	"pairs-trading-backend/internal/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := database.Init()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	server := http.NewServer(cfg, db)

	log.Printf("Starting server on :%s", cfg.Port)
	if err := server.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
