package repo

import (
	"errors"
	"telemafia/internal/scenario/entity"
)

// AssignmentRepository defines the interface for storing role assignments
type AssignmentRepository interface {
	StoreAssignment(roomID string, assignments map[string]entity.Role) error
	GetAssignment(roomID string) (map[string]entity.Role, error)
	StoreRoomScenario(roomID string, scenarioName string) error
	GetRoomScenario(roomID string) (string, error)
}

// InMemoryAssignmentRepository provides an in-memory implementation of AssignmentRepository
type InMemoryAssignmentRepository struct {
	data           map[string]map[string]entity.Role
	roomToScenario map[string]string // Maps roomID to scenarioName
}

// NewInMemoryAssignmentRepository creates a new in-memory assignment repository
func NewInMemoryAssignmentRepository() *InMemoryAssignmentRepository {
	return &InMemoryAssignmentRepository{
		data:           make(map[string]map[string]entity.Role),
		roomToScenario: make(map[string]string),
	}
}

// StoreAssignment stores the role assignments for a room
func (r *InMemoryAssignmentRepository) StoreAssignment(roomID string, assignments map[string]entity.Role) error {
	r.data[roomID] = assignments
	return nil
}

// GetAssignment retrieves the role assignments for a room
func (r *InMemoryAssignmentRepository) GetAssignment(roomID string) (map[string]entity.Role, error) {
	assignments, exists := r.data[roomID]
	if !exists {
		return nil, errors.New("assignments not found")
	}
	return assignments, nil
}

// StoreRoomScenario stores the scenario name for a room
func (r *InMemoryAssignmentRepository) StoreRoomScenario(roomID string, scenarioName string) error {
	r.roomToScenario[roomID] = scenarioName
	return nil
}

// GetRoomScenario retrieves the scenario name for a room
func (r *InMemoryAssignmentRepository) GetRoomScenario(roomID string) (string, error) {
	scenarioName, exists := r.roomToScenario[roomID]
	if !exists {
		return "", errors.New("no scenario assigned to room")
	}
	return scenarioName, nil
}
