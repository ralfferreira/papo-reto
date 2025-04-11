package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ralfferreira/papo-reto/internal/config"
	"github.com/ralfferreira/papo-reto/internal/repository"
	"github.com/ralfferreira/papo-reto/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := repository.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate database schema
	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Create and start server
	srv := server.NewServer(cfg, db)
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %s", cfg.Server.Port)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := srv.Shutdown(); err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	}

	log.Println("Server exited properly")
}
