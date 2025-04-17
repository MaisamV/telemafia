package usecase

import (
	"context"
	"telemafia/internal/entity"
	"telemafia/pkg/event"
	"time"
)

// JoinRoomCommand represents the command to join a room
type JoinRoomCommand struct {
	Requester entity.User
	RoomID    entity.RoomID
}

// JoinRoomHandler handles room joining
type JoinRoomHandler struct {
	roomRepo       Repository
	eventPublisher event.Publisher
}

// NewJoinRoomHandler creates a new JoinRoomHandler
func NewJoinRoomHandler(repo Repository, publisher event.Publisher) *JoinRoomHandler {
	return &JoinRoomHandler{
		roomRepo:       repo,
		eventPublisher: publisher,
	}
}

// Handle processes the join room command
func (h *JoinRoomHandler) Handle(ctx context.Context, cmd JoinRoomCommand) error {
	// Get the room first to ensure it exists
	if _, err := h.roomRepo.GetRoomByID(cmd.RoomID); err != nil {
		return err
	}

	// Add player to room
	if err := h.roomRepo.AddPlayerToRoom(cmd.RoomID, &cmd.Requester); err != nil {
		return err
	}

	// Publish domain event
	event := entity.PlayerJoinedEvent{
		RoomID:   cmd.RoomID,
		PlayerID: cmd.Requester.ID,
		JoinedAt: time.Now(),
	}

	if err := h.eventPublisher.Publish(event); err != nil {
		// Log error but don't fail the operation
		// Consider using a retry mechanism for event publishing
	}

	return nil
}
