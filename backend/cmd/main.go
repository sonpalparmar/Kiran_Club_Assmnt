// cmd/server/main.go
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"kirana-club-assignment/backend/config"
	"kirana-club-assignment/backend/internal/api"
	"kirana-club-assignment/backend/internal/processor"
	"kirana-club-assignment/backend/internal/repository"
	"kirana-club-assignment/backend/internal/service"
)

func main() {
	// Initialize random seed for simulated processing time
	rand.Seed(time.Now().UnixNano())

	// Load configuration
	cfg := config.Load()

	// Initialize repositories
	jobRepo := repository.NewInMemoryJobRepository()
	storeRepo := repository.NewInMemoryStoreRepository()

	// Load stores from file
	if err := storeRepo.LoadStoresFromFile(cfg.StoreMasterFile); err != nil {
		log.Fatalf("Failed to load store data: %v", err)
	}

	// Initialize image processor
	imageProcessor := processor.NewImageProcessor()

	// Initialize job service
	jobService := service.NewJobService(jobRepo, storeRepo, imageProcessor)

	// Initialize API handler
	handler := api.NewHandler(jobService)

	// Setup routes
	router := api.SetupRoutes(handler)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
