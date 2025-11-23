package repository

import (
	"delivery-state-manager/internal/models"
	"delivery-state-manager/pkg/errs"
	"sync"
)

// Repository defines the interface for data access operations
type Repository interface {
	// Driver operations
	CreateOrUpdateDriver(driver *models.Driver)
	GetDriver(id string) (*models.Driver, error)
	GetAllDrivers() []*models.Driver
	UpdateDriverStatus(id string, status models.DriverStatus) error
	GetAvailableDrivers() []*models.Driver

	// Order operations
	CreateOrder(order *models.Order)
	GetOrder(id string) (*models.Order, error)
	GetAllOrders() []*models.Order
	UpdateOrderStatus(id string, status models.OrderStatus) error
	GetPendingOrders() []*models.Order

	// Assignment operations
	AssignOrderToDriver(orderID, driverID string) error

	// Debug operations
	GetSnapshot() models.StateSnapshot
}

// StateManager manages all drivers and orders with thread-safe access
type StateManager struct {
	drivers map[string]*models.Driver
	orders  map[string]*models.Order
	mu      sync.RWMutex
}

// NewStateManager creates a new StateManager instance
func NewStateManager() Repository {
	return &StateManager{
		drivers: make(map[string]*models.Driver),
		orders:  make(map[string]*models.Order),
	}
}

// CreateOrUpdateDriver creates a new driver or updates an existing one
func (sm *StateManager) CreateOrUpdateDriver(driver *models.Driver) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	driver.UpdatedAt = models.GetCurrentTimestamp()
	sm.drivers[driver.ID] = driver
}

// GetDriver retrieves a driver by ID
func (sm *StateManager) GetDriver(id string) (*models.Driver, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	driver, ok := sm.drivers[id]
	if !ok {
		return nil, errs.ErrDriverNotFound
	}

	// Return a copy to prevent external mutation
	driverCopy := *driver
	return &driverCopy, nil
}

// GetAllDrivers returns all drivers
func (sm *StateManager) GetAllDrivers() []*models.Driver {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	drivers := make([]*models.Driver, 0, len(sm.drivers))
	for _, driver := range sm.drivers {
		driverCopy := *driver
		drivers = append(drivers, &driverCopy)
	}
	return drivers
}

// UpdateDriverStatus updates the status of a driver
func (sm *StateManager) UpdateDriverStatus(id string, status models.DriverStatus) error {
	if !models.IsValidDriverStatus(status) {
		return errs.ErrInvalidStatusUpdate
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	driver, ok := sm.drivers[id]
	if !ok {
		return errs.ErrDriverNotFound
	}

	driver.Status = status
	driver.UpdatedAt = models.GetCurrentTimestamp()
	return nil
}

// CreateOrder creates a new order with pending status
func (sm *StateManager) CreateOrder(order *models.Order) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	now := models.GetCurrentTimestamp()
	order.Status = models.OrderPending
	order.CreatedAt = now
	order.UpdatedAt = now
	order.DriverID = ""

	sm.orders[order.ID] = order
}

// GetOrder retrieves an order by ID
func (sm *StateManager) GetOrder(id string) (*models.Order, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	order, ok := sm.orders[id]
	if !ok {
		return nil, errs.ErrOrderNotFound
	}

	// Return a copy to prevent external mutation
	orderCopy := *order
	return &orderCopy, nil
}

// GetAllOrders returns all orders
func (sm *StateManager) GetAllOrders() []*models.Order {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	orders := make([]*models.Order, 0, len(sm.orders))
	for _, order := range sm.orders {
		orderCopy := *order
		orders = append(orders, &orderCopy)
	}
	return orders
}

// UpdateOrderStatus updates the status of an order with validation
func (sm *StateManager) UpdateOrderStatus(id string, status models.OrderStatus) error {
	if !models.IsValidOrderStatus(status) {
		return errs.ErrInvalidStatusUpdate
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	order, ok := sm.orders[id]
	if !ok {
		return errs.ErrOrderNotFound
	}

	// Validate state transition
	if !models.CanTransitionOrderStatus(order.Status, status) {
		return errs.ErrInvalidTransition
	}

	order.Status = status
	order.UpdatedAt = models.GetCurrentTimestamp()
	return nil
}

// GetPendingOrders returns all orders with pending status
func (sm *StateManager) GetPendingOrders() []*models.Order {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	pending := make([]*models.Order, 0)
	for _, order := range sm.orders {
		if order.Status == models.OrderPending {
			orderCopy := *order
			pending = append(pending, &orderCopy)
		}
	}
	return pending
}

// GetAvailableDrivers returns all drivers with available status
func (sm *StateManager) GetAvailableDrivers() []*models.Driver {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	available := make([]*models.Driver, 0)
	for _, driver := range sm.drivers {
		if driver.Status == models.DriverAvailable {
			driverCopy := *driver
			available = append(available, &driverCopy)
		}
	}
	return available
}

// AssignOrderToDriver atomically assigns an order to a driver
func (sm *StateManager) AssignOrderToDriver(orderID, driverID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	order, ok := sm.orders[orderID]
	if !ok {
		return errs.ErrOrderNotFound
	}

	driver, ok := sm.drivers[driverID]
	if !ok {
		return errs.ErrDriverNotFound
	}

	// Validate order status
	if order.Status != models.OrderPending {
		return errs.ErrOrderAlreadyAssigned
	}

	// Validate driver status
	if driver.Status != models.DriverAvailable {
		return errs.ErrDriverNotAvailable
	}

	// Perform atomic assignment
	order.Status = models.OrderAssigned
	order.DriverID = driverID
	order.UpdatedAt = models.GetCurrentTimestamp()

	driver.Status = models.DriverBusy
	driver.UpdatedAt = models.GetCurrentTimestamp()

	return nil
}

// GetSnapshot returns a complete snapshot of the current state
func (sm *StateManager) GetSnapshot() models.StateSnapshot {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	snapshot := models.StateSnapshot{
		Drivers:   make(map[string]*models.Driver),
		Orders:    make(map[string]*models.Order),
		Timestamp: models.GetCurrentTimestamp(),
	}

	for id, driver := range sm.drivers {
		driverCopy := *driver
		snapshot.Drivers[id] = &driverCopy
	}

	for id, order := range sm.orders {
		orderCopy := *order
		snapshot.Orders[id] = &orderCopy
	}

	return snapshot
}
