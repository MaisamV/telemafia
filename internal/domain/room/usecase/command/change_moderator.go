package command

import (
	"context"
	"fmt"
	roomEntity "telemafia/internal/domain/room/entity"
	roomPort "telemafia/internal/domain/room/port"
	sharedEntity "telemafia/internal/shared/entity"
	// sharedEvent "telemafia/internal/shared/event" // No event needed for now
)

// ChangeModeratorCommand represents the command to change a room's moderator
type ChangeModeratorCommand struct {
	Requester    *sharedEntity.User // User initiating the change
	RoomID       roomEntity.RoomID
	NewModerator *sharedEntity.User // User to become the new moderator
}

// ChangeModeratorHandler handles changing the room moderator
type ChangeModeratorHandler struct {
	roomRepo roomPort.RoomRepository // Need full repo to get/update
	// eventPublisher sharedEvent.Publisher // If events are added later
}

// NewChangeModeratorHandler creates a new ChangeModeratorHandler
func NewChangeModeratorHandler(repo roomPort.RoomRepository) *ChangeModeratorHandler {
	return &ChangeModeratorHandler{
		roomRepo: repo,
	}
}

// Handle processes the change moderator command
func (h *ChangeModeratorHandler) Handle(ctx context.Context, cmd ChangeModeratorCommand) error {
	// --- Basic Validation ---
	if cmd.Requester == nil {
		return fmt.Errorf("change moderator: requester cannot be nil")
	}
	if cmd.NewModerator == nil {
		return fmt.Errorf("change moderator: new moderator cannot be nil")
	}

	// Fetch the room
	room, err := h.roomRepo.GetRoomByID(cmd.RoomID)
	if err != nil {
		return fmt.Errorf("change moderator: could not find room %s: %w", cmd.RoomID, err)
	}

	// --- Permission Check ---
	// Allow if requester is global admin OR the current moderator of this specific room
	isCurrentModerator := room.Moderator != nil && room.Moderator.ID == cmd.Requester.ID
	if !cmd.Requester.Admin && !isCurrentModerator {
		return fmt.Errorf("change moderator: permission denied (requires admin or current room moderator)")
	}

	// Ensure the new moderator is not the current moderator
	if room.Moderator != nil && room.Moderator.ID == cmd.NewModerator.ID {
		return fmt.Errorf("change moderator: user %s is already the moderator", cmd.NewModerator.GetProfileLink())
	}

	// Use the entity method to set the new moderator
	if err := room.SetModerator(cmd.NewModerator); err != nil {
		return fmt.Errorf("change moderator: failed to set new moderator: %w", err)
	}

	// Update the room in the repository
	if err := h.roomRepo.UpdateRoom(room); err != nil {
		return fmt.Errorf("change moderator: failed to save room updates: %w", err)
	}

	// TODO: Optionally publish a RoomModeratorChangedEvent

	return nil
}
