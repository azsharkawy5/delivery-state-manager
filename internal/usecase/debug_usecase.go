package usecase

import (
	"delivery-state-manager/internal/models"
)

// DebugRepository defines the interface for debug operations
type DebugRepository interface {
	GetSnapshot() models.StateSnapshot
}

// DebugUseCase handles debug-related use cases
type DebugUseCase struct {
	repo DebugRepository
}

// NewDebugUseCase creates a new DebugUseCase instance
func NewDebugUseCase(repo DebugRepository) *DebugUseCase {
	return &DebugUseCase{
		repo: repo,
	}
}

// GetSnapshot returns a complete snapshot of the current state
func (uc *DebugUseCase) GetSnapshot() models.StateSnapshot {
	return uc.repo.GetSnapshot()
}
