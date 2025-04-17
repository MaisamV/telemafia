package usecase

import (
	"context"
	"errors"
	"telemafia/internal/entity"
	"telemafia/pkg/event"
	"time"
)

// KickUserCommand represents the command to kick a user from a room
type KickUserCommand struct {
	Requester entity.User
	RoomID    entity.RoomID
	PlayerID  entity.UserID
}

// KickUserHandler handles kicking a user from a room
type KickUserHandler struct {
	roomRepo       Repository
	eventPublisher event.Publisher
}

// NewKickUserHandler creates a new KickUserHandler
func NewKickUserHandler(repo Repository, publisher event.Publisher) *KickUserHandler {
	return &KickUserHandler{
		roomRepo:       repo,
		eventPublisher: publisher,
	}
}

// Handle processes the kick user command
func (h *KickUserHandler) Handle(ctx context.Context, cmd KickUserCommand) error {
	if !cmd.Requester.Admin {
		return errors.New("admin privilege required")
	}
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
