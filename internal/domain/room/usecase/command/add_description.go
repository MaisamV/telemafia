package command

import (
	"context"
	// roomEntity "telemafia/internal/room/entity" // Room is needed here
	roomEntity "telemafia/internal/domain/room/entity"
	// roomPort "telemafia/internal/room/port" // Room repo is needed
	roomPort "telemafia/internal/domain/room/port"
	// scenarioEntity "telemafia/internal/scenario/entity"
	// scenarioPort "telemafia/internal/scenario/port"
	"errors"
	"fmt"
	sharedEntity "telemafia/internal/shared/entity" // Added for User
)

// AddDescriptionCommand represents the command to add a description to a room
// NOTE: This seems more related to the Room domain than Scenario based on dependencies.
// Consider moving this to internal/room/usecase/
type AddDescriptionCommand struct {
	Requester       sharedEntity.User // The user initiating the action
	Room            *roomEntity.Room  // Pass the Room object directly
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
	// --- Permission Check ---
	if !cmd.Requester.Admin {
		return errors.New("add description: admin privilege required")
	}

	if cmd.Room == nil {
		return errors.New("cannot add description to nil room")
	}

	// Modify the passed-in room object
	cmd.Room.SetDescription(cmd.DescriptionName, cmd.Text)

	// Persist the changes using UpdateRoom
	if err := h.roomRepo.UpdateRoom(cmd.Room); err != nil {
		return fmt.Errorf("failed to update room after adding description: %w", err)
	}

	// h.roomRepo.RaiseChangeFlag() // No longer needed, UpdateRoom handles the flag
	return nil
}
