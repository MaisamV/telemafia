package command

import (
	"context"
	"telemafia/internal/room/entity"
	"telemafia/internal/room/repo"
	userEntity "telemafia/internal/user/entity"
	"telemafia/pkg/event"
	"time"
)

// LeaveRoomCommand represents the command to leave a room
type LeaveRoomCommand struct {
	Requester userEntity.User
	RoomID    entity.RoomID
}

// LeaveRoomHandler handles room leaving
type LeaveRoomHandler struct {
	roomRepo       repo.Repository
	eventPublisher event.Publisher
}

// NewLeaveRoomHandler creates a new LeaveRoomHandler
func NewLeaveRoomHandler(repo repo.Repository, publisher event.Publisher) *LeaveRoomHandler {
	return &LeaveRoomHandler{
		roomRepo:       repo,
		eventPublisher: publisher,
	}
}

// Handle processes the leave room command
func (h *LeaveRoomHandler) Handle(ctx context.Context, cmd LeaveRoomCommand) error {
	// Get the room first to ensure it exists
	if _, err := h.roomRepo.GetRoomByID(cmd.RoomID); err != nil {
		return err
	}

	// Remove player from room
	if err := h.roomRepo.RemovePlayerFromRoom(cmd.RoomID, cmd.Requester.ID); err != nil {
		return err
	}

	// Publish domain event
	event := entity.PlayerLeftEvent{
		RoomID:   cmd.RoomID,
		PlayerID: cmd.Requester.ID,
		LeftAt:   time.Now(),
	}

	if err := h.eventPublisher.Publish(event); err != nil {
		// Log error but don't fail the operation
		// Consider using a retry mechanism for event publishing
	}

	return nil
}
