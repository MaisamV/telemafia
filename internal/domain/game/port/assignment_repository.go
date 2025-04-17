package port

import (
	scenarioEntity "telemafia/internal/domain/scenario/entity"
	// Assuming roomID is a string here, might need roomEntity import if it's RoomID type
)

// AssignmentRepository defines the interface for storing role assignments
// TODO: Review if this interface is still needed or if functionality should be moved
//
//	into GameRepository (for Assignments) and RoomRepository (for RoomScenario).
type AssignmentRepository interface {
	// StoreAssignment associates role assignments with a room ID (likely should be GameID)
	StoreAssignment(roomID string, assignments map[string]scenarioEntity.Role) error
	// GetAssignment retrieves role assignments for a room ID (likely should be GameID)
	GetAssignment(roomID string) (map[string]scenarioEntity.Role, error)
	// StoreRoomScenario associates a scenario name with a room ID
	StoreRoomScenario(roomID string, scenarioName string) error
	// GetRoomScenario retrieves the scenario name for a room ID
	GetRoomScenario(roomID string) (string, error)
}
