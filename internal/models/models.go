package models

import "time"

// Location represents a geographic coordinate
type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// DriverStatus represents the current status of a driver
type DriverStatus string

const (
	DriverAvailable DriverStatus = "available"
	DriverBusy      DriverStatus = "busy"
	DriverOffline   DriverStatus = "offline"
)

// Driver represents a delivery driver
type Driver struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Status    DriverStatus `json:"status"`
	Location  Location     `json:"location"`
	UpdatedAt int64        `json:"updated_at"`
}

// OrderStatus represents the current status of an order
type OrderStatus string

const (
	OrderPending   OrderStatus = "pending"
	OrderAssigned  OrderStatus = "assigned"
	OrderPickedUp  OrderStatus = "picked_up"
	OrderDelivered OrderStatus = "delivered"
	OrderCanceled  OrderStatus = "canceled"
)

// Order represents a customer order
type Order struct {
	ID        string      `json:"id"`
	Customer  string      `json:"customer"`
	Pickup    Location    `json:"pickup"`
	Dropoff   Location    `json:"dropoff"`
	Status    OrderStatus `json:"status"`
	DriverID  string      `json:"driver_id,omitempty"`
	CreatedAt int64       `json:"created_at"`
	UpdatedAt int64       `json:"updated_at"`
}

// StateSnapshot represents a complete snapshot of the system state
type StateSnapshot struct {
	Drivers   map[string]*Driver `json:"drivers"`
	Orders    map[string]*Order  `json:"orders"`
	Timestamp int64              `json:"timestamp"`
}

// ===== Utility Functions =====

// IsValidDriverStatus checks if a driver status is valid
func IsValidDriverStatus(status DriverStatus) bool {
	switch status {
	case DriverAvailable, DriverBusy, DriverOffline:
		return true
	}
	return false
}

// IsValidOrderStatus checks if an order status is valid
func IsValidOrderStatus(status OrderStatus) bool {
	switch status {
	case OrderPending, OrderAssigned, OrderPickedUp, OrderDelivered, OrderCanceled:
		return true
	}
	return false
}

// CanTransitionOrderStatus checks if an order status transition is valid
func CanTransitionOrderStatus(from, to OrderStatus) bool {
	// Define valid state transitions
	validTransitions := map[OrderStatus][]OrderStatus{
		OrderPending:   {OrderAssigned, OrderCanceled},
		OrderAssigned:  {OrderPickedUp, OrderCanceled},
		OrderPickedUp:  {OrderDelivered, OrderCanceled},
		OrderDelivered: {},
		OrderCanceled:  {},
	}

	allowedStates, ok := validTransitions[from]
	if !ok {
		return false
	}

	for _, allowed := range allowedStates {
		if allowed == to {
			return true
		}
	}
	return false
}

// GetCurrentTimestamp returns the current Unix timestamp
func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}
