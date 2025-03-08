package command

import (
	"context"
	"telemafia/internal/room/entity"
	"telemafia/internal/room/repo"
	userEntity "telemafia/internal/user/entity"
	"telemafia/pkg/event"
	"time"
)

// KickUserCommand represents the command to kick a user from a room
type KickUserCommand struct {
	RoomID   entity.RoomID
	PlayerID userEntity.UserID
}

// KickUserHandler handles kicking a user from a room
type KickUserHandler struct {
	roomRepo       repo.Repository
	eventPublisher event.Publisher
}

// NewKickUserHandler creates a new KickUserHandler
func NewKickUserHandler(repo repo.Repository, publisher event.Publisher) *KickUserHandler {
	return &KickUserHandler{
		roomRepo:       repo,
		eventPublisher: publisher,
	}
}

// Handle processes the kick user command
func (h *KickUserHandler) Handle(ctx context.Context, cmd KickUserCommand) error {
	// Get the room first to ensure it exists
	if _, err := h.roomRepo.GetRoomByID(cmd.RoomID); err != nil {
		return err
	}

	// Remove player from room
	if err := h.roomRepo.RemovePlayerFromRoom(cmd.RoomID, cmd.PlayerID); err != nil {
		return err
	}

	// Publish domain event
	event := entity.PlayerKickedEvent{
		RoomID:   cmd.RoomID,
		PlayerID: cmd.PlayerID,
		KickedAt: time.Now(),
	}

	if err := h.eventPublisher.Publish(event); err != nil {
		// Log error but don't fail the operation
		// Consider using a retry mechanism for event publishing
	}

	return nil
}
