package usecase

import (
	"delivery-state-manager/internal/models"
	"delivery-state-manager/pkg/errs"
)

// DriverRepository defines the interface for driver operations
type DriverRepository interface {
	CreateOrUpdateDriver(driver *models.Driver)
	GetDriver(id string) (*models.Driver, error)
	GetAllDrivers() []*models.Driver
	UpdateDriverStatus(id string, status models.DriverStatus) error
}

// DriverUseCase handles driver-related use cases
type DriverUseCase struct {
	repo DriverRepository
}

// NewDriverUseCase creates a new DriverUseCase instance
func NewDriverUseCase(repo DriverRepository) *DriverUseCase {
	return &DriverUseCase{
		repo: repo,
	}
}

// CreateOrUpdateDriver creates or updates a driver
func (uc *DriverUseCase) CreateOrUpdateDriver(driver *models.Driver) error {
	// Validate required fields
	if driver.ID == "" || driver.Name == "" {
		return errs.ErrMissingRequiredField
	}

	// Validate status if provided
	if driver.Status != "" && !models.IsValidDriverStatus(driver.Status) {
		return errs.ErrInvalidStatusUpdate
	}

	// Set default status if not provided
	if driver.Status == "" {
		driver.Status = models.DriverAvailable
	}

	uc.repo.CreateOrUpdateDriver(driver)
	return nil
}

// GetDriver retrieves a driver by ID
func (uc *DriverUseCase) GetDriver(id string) (*models.Driver, error) {
	return uc.repo.GetDriver(id)
}

// GetAllDrivers returns all drivers
func (uc *DriverUseCase) GetAllDrivers() []*models.Driver {
	return uc.repo.GetAllDrivers()
}

// UpdateDriverStatus updates the status of a driver
func (uc *DriverUseCase) UpdateDriverStatus(id string, status models.DriverStatus) error {
	return uc.repo.UpdateDriverStatus(id, status)
}
