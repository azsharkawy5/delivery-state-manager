package service

import (
	"delivery-state-manager/internal/models"
	"log"
	"time"
)

// MatcherRepository defines the interface for the matching repository
type MatcherRepository interface {
	AssignOrderToDriver(orderID, driverID string) error
	GetAvailableDrivers() []*models.Driver
	GetPendingOrders() []*models.Order
}

// Matcher handles order-to-driver matching
type Matcher struct {
	repo MatcherRepository
}

// NewMatcher creates a new Matcher instance
func NewMatcher(repo MatcherRepository) *Matcher {
	return &Matcher{
		repo: repo,
	}
}

// StartMatcher runs the background matching engine
func (m *Matcher) StartMatcher(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("Matcher started with interval: %v", interval)

	for range ticker.C {
		m.MatchOrders()
	}
}

// MatchOrders performs the actual matching logic
func (m *Matcher) MatchOrders() {
	pendingOrders := m.repo.GetPendingOrders()
	availableDrivers := m.repo.GetAvailableDrivers()

	if len(pendingOrders) == 0 {
		return
	}

	if len(availableDrivers) == 0 {
		log.Printf("No available drivers for %d pending orders", len(pendingOrders))
		return
	}

	matched := 0

	// Simple first-come-first-served matching
	for i, order := range pendingOrders {
		if i >= len(availableDrivers) {
			break
		}

		driver := availableDrivers[i]

		err := m.repo.AssignOrderToDriver(order.ID, driver.ID)
		if err != nil {
			log.Printf("Failed to assign order %s to driver %s: %v", order.ID, driver.ID, err)
			continue
		}

		log.Printf("Matched order %s to driver %s", order.ID, driver.ID)
		matched++
	}

	if matched > 0 {
		log.Printf("Matcher completed: %d orders assigned to drivers", matched)
	}
}
