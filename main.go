package main

import (
	"delivery-state-manager/config"
	"delivery-state-manager/internal/handler"
	"delivery-state-manager/internal/repository"
	"delivery-state-manager/internal/service"
	"delivery-state-manager/internal/usecase"

	"log"
)

func main() {
	log.Println("Starting Delivery State Manager...")
	// Load config
	config := config.LoadConfig()

	// Initialize repository layer
	repo := repository.NewStateManager()

	// Initialize service layer
	matcherService := service.NewMatcher(repo)

	// Initialize use case layer
	driverUC := usecase.NewDriverUseCase(repo)
	orderUC := usecase.NewOrderUseCase(repo)
	debugUC := usecase.NewDebugUseCase(repo)

	// Initialize handler layer
	h := handler.NewHandler(driverUC, orderUC, debugUC)

	// Start background matcher
	go matcherService.StartMatcher(config.MatcherInterval)

	// Setup HTTP router
	router := h.SetupRouter()

	// Start HTTP server
	log.Printf("Server listening on %s", config.ServerPort)
	if err := router.Run(config.ServerPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
