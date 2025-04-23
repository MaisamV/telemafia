package command

import (
	"context"
	"errors"
	roomEntity "telemafia/internal/domain/room/entity"
	roomPort "telemafia/internal/domain/room/port"
	sharedEntity "telemafia/internal/shared/entity"
	sharedEvent "telemafia/internal/shared/event"
	"time"
)

// KickUserCommand represents the command to kick a user from a room
type KickUserCommand struct {
	Requester sharedEntity.User   // The user initiating the kick
	RoomID    roomEntity.RoomID   // Use imported RoomID type
	PlayerID  sharedEntity.UserID // Use imported UserID type
}

// KickUserHandler handles kicking a user from a room
type KickUserHandler struct {
	roomRepo       roomPort.RoomRepository // Use imported Repository interface
	eventPublisher sharedEvent.Publisher
}

// NewKickUserHandler creates a new KickUserHandler
func NewKickUserHandler(repo roomPort.RoomRepository, publisher sharedEvent.Publisher) *KickUserHandler {
	return &KickUserHandler{
		roomRepo:       repo,
		eventPublisher: publisher,
	}
}

// Handle processes the kick user command
func (h *KickUserHandler) Handle(ctx context.Context, cmd KickUserCommand) error {
	// --- Permission Check ---
	if !cmd.Requester.Admin { // Assuming Admin field exists on sharedEntity.User
		return errors.New("kick user: admin privilege required") // More specific error
	}

	// Check if the room exists (optional, RemovePlayerFromRoom likely handles not found)
	_, err := h.roomRepo.GetRoomByID(cmd.RoomID)
	if err != nil {
		return err
	}

	// Remove player from room
	if err := h.roomRepo.RemovePlayerFromRoom(cmd.RoomID, cmd.PlayerID); err != nil {
		return err // Propagates ErrPlayerNotInRoom etc.
	}

	// Publish domain event
	evt := sharedEvent.PlayerKickedEvent{ // Use imported event type
		RoomID:   cmd.RoomID,
		PlayerID: cmd.PlayerID,
		KickedAt: time.Now(),
	}

	if err := h.eventPublisher.Publish(evt); err != nil {
		// Log error but don't fail the operation
		// log.Printf("Failed to publish PlayerKickedEvent: %v", err)
	}

	return nil
}
