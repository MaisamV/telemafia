package command

import (
	"context"
	"telemafia/internal/room/entity"
	"telemafia/internal/room/repo"
)

// AssignScenarioCommand represents a command to assign a scenario to a room
type AssignScenarioCommand struct {
	RoomID       entity.RoomID
	ScenarioName string
}

// AssignScenarioHandler handles the assign scenario command
type AssignScenarioHandler struct {
	roomRepo repo.Repository
}

// NewAssignScenarioHandler creates a new assign scenario handler
func NewAssignScenarioHandler(roomRepo repo.Repository) *AssignScenarioHandler {
	return &AssignScenarioHandler{
		roomRepo: roomRepo,
	}
}

// Handle assigns a scenario to a room
func (h *AssignScenarioHandler) Handle(ctx context.Context, cmd AssignScenarioCommand) error {
	// First, verify the room exists
	room, err := h.roomRepo.GetRoomByID(cmd.RoomID)
	if err != nil {
		return err
	}

	// Assign the scenario to the room
	err = h.roomRepo.AssignScenarioToRoom(cmd.RoomID, cmd.ScenarioName)
	if err != nil {
		return err
	}

	// Update the room entity
	room.AssignScenario(cmd.ScenarioName)
	return nil
}
