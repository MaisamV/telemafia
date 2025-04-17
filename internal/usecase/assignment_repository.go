package usecase

import (
	// "errors" // Removed unused import
	"telemafia/internal/entity"
)

// AssignmentRepository defines the interface for storing role assignments
// TODO: Review if this interface is still needed or if functionality is covered by GameRepository/RoomRepository
type AssignmentRepository interface {
	StoreAssignment(roomID string, assignments map[string]entity.Role) error
	GetAssignment(roomID string) (map[string]entity.Role, error)
	StoreRoomScenario(roomID string, scenarioName string) error
	GetRoomScenario(roomID string) (string, error)
}
