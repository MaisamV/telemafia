package command

import (
	"context"
	roomEntity "telemafia/internal/domain/room/entity"
	roomPort "telemafia/internal/domain/room/port"
	sharedEntity "telemafia/internal/shared/entity"
	sharedEvent "telemafia/internal/shared/event"
	"time"
)

// LeaveRoomCommand represents the command to leave a room
type LeaveRoomCommand struct {
	Requester sharedEntity.User // Use imported User type
	RoomID    roomEntity.RoomID // Use imported RoomID type
}

// LeaveRoomHandler handles room leaving
type LeaveRoomHandler struct {
	roomRepo       roomPort.RoomRepository // Use imported Repository interface
	eventPublisher sharedEvent.Publisher
}

// NewLeaveRoomHandler creates a new LeaveRoomHandler
func NewLeaveRoomHandler(repo roomPort.RoomRepository, publisher sharedEvent.Publisher) *LeaveRoomHandler {
	return &LeaveRoomHandler{
		roomRepo:       repo,
		eventPublisher: publisher,
	}
}

// Handle processes the leave room command
func (h *LeaveRoomHandler) Handle(ctx context.Context, cmd LeaveRoomCommand) error {
	// Check if the room exists (optional, RemovePlayerFromRoom might handle this)
	_, err := h.roomRepo.GetRoomByID(cmd.RoomID)
	if err != nil {
		return err
	}

	// Remove player from room
	if err := h.roomRepo.RemovePlayerFromRoom(cmd.RoomID, cmd.Requester.ID); err != nil {
		return err // Propagates ErrPlayerNotInRoom etc.
	}

	// Publish domain event
	evt := sharedEvent.PlayerLeftEvent{ // Use imported event type
		RoomID:   cmd.RoomID,
		PlayerID: cmd.Requester.ID,
		LeftAt:   time.Now(),
	}

	if err := h.eventPublisher.Publish(evt); err != nil {
		// Log error but don't fail the operation
		// log.Printf("Failed to publish PlayerLeftEvent: %v", err)
	}

	return nil
}
