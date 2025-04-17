package command

import (
	"context"
	// roomEntity "telemafia/internal/room/entity" // Room is needed here
	roomEntity "telemafia/internal/domain/room/entity"
	// roomPort "telemafia/internal/room/port" // Room repo is needed
	roomPort "telemafia/internal/domain/room/port"
	// scenarioEntity "telemafia/internal/scenario/entity"
	// scenarioPort "telemafia/internal/scenario/port"
)

// AddDescriptionCommand represents the command to add a description to a room
// NOTE: This seems more related to the Room domain than Scenario based on dependencies.
// Consider moving this to internal/room/usecase/
type AddDescriptionCommand struct {
	RoomID          roomEntity.RoomID // Use imported RoomID type
	DescriptionName string
	Text            string
}

// AddDescriptionHandler handles adding description to a room
type AddDescriptionHandler struct {
	roomRepo roomPort.RoomRepository // Use imported RoomRepository interface
}

// NewAddDescriptionHandler creates a new AddDescriptionHandler
func NewAddDescriptionHandler(repo roomPort.RoomRepository) *AddDescriptionHandler {
	return &AddDescriptionHandler{roomRepo: repo}
}

// Handle processes the add description command
func (h *AddDescriptionHandler) Handle(ctx context.Context, cmd AddDescriptionCommand) error {
	room, err := h.roomRepo.GetRoomByID(cmd.RoomID) // Get the room
	if err != nil {
		return err // Propagates ErrRoomNotFound etc.
	}

	// Assuming the entity method handles updating the map correctly.
	// The persistence of this change might require an UpdateRoom method in the repo.
	room.SetDescription(cmd.DescriptionName, cmd.Text)

	// TODO: Need to call an UpdateRoom method on the repository here?
	// Example: return h.roomRepo.UpdateRoom(room)
	// For now, assume the in-memory object reference is sufficient (dangerous for real persistence).
	h.roomRepo.RaiseChangeFlag() // Mark as changed for potential in-memory save
	return nil
}
